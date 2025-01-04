package services

import (
	"context"
	"fmt"
	"github.com/massimo-ua/quill/internal/domain"
	"github.com/massimo-ua/quill/internal/domain/ports"
)

type MessageHandler interface {
	Handle(ctx context.Context, msg *domain.Message) error
}

type SuggestionsHandler interface {
	MessageHandler
	HandleWithAnalysis(ctx context.Context, msg *domain.Message, analysis *domain.MessageAnalysisResult) error
}

type baseHandler struct {
	docService   *DocumentationService
	chatProvider ports.ChatAccessProvider
}

type ideaHandler struct {
	baseHandler
}

type decisionHandler struct {
	baseHandler
}

type statusHandler struct {
	baseHandler
}

type unknownHandler struct {
	baseHandler
}

func (h *baseHandler) createDocumentation(ctx context.Context, msg *domain.Message) error {
	return h.docService.CreateDocumentation(
		ctx,
		msg.Type(),
		msg.Category(),
		msg.Content().Text(),
		msg.References(),
	)
}

func (h *ideaHandler) Handle(ctx context.Context, msg *domain.Message) error {
	if err := h.createDocumentation(ctx, msg); err != nil {
		return fmt.Errorf("failed to create idea documentation: %w", err)
	}

	reply := fmt.Sprintf("📝 Captured idea in category: %s", msg.Category())
	if msg.HasReferences() {
		reply += fmt.Sprintf("\n🔗 Linked to %d related items", len(msg.References()))
	}

	return h.chatProvider.ReplyToMessage(ctx, msg.ID().String(), reply)
}

func (h *decisionHandler) Handle(ctx context.Context, msg *domain.Message) error {
	if err := h.createDocumentation(ctx, msg); err != nil {
		return fmt.Errorf("failed to create decision documentation: %w", err)
	}

	reply := fmt.Sprintf("✅ Recorded decision in category: %s", msg.Category())
	if msg.HasReferences() {
		reply += fmt.Sprintf("\n🔗 Linked to %d related items", len(msg.References()))
	}

	return h.chatProvider.ReplyToMessage(ctx, msg.ID().String(), reply)
}

func (h *statusHandler) Handle(ctx context.Context, msg *domain.Message) error {
	if err := h.createDocumentation(ctx, msg); err != nil {
		return fmt.Errorf("failed to create status documentation: %w", err)
	}

	reply := fmt.Sprintf("📊 Logged status update in category: %s", msg.Category())
	if msg.HasReferences() {
		reply += fmt.Sprintf("\n🔗 Linked to %d related items", len(msg.References()))
	}

	return h.chatProvider.ReplyToMessage(ctx, msg.ID().String(), reply)
}

func (h *unknownHandler) Handle(ctx context.Context, msg *domain.Message) error {
	return nil
}

func (h *unknownHandler) HandleWithAnalysis(ctx context.Context, msg *domain.Message, analysis *domain.MessageAnalysisResult) error {
	if !analysis.IsHighConfidence() || !analysis.HasSuggestedTags() {
		return nil
	}

	suggestion := "💡 I noticed this might be relevant. Consider adding these tags:\n"
	for _, tag := range analysis.SuggestedTags() {
		suggestion += fmt.Sprintf("- #%s\n", tag)
	}

	return h.chatProvider.ReplyToMessage(ctx, msg.ID().String(), suggestion)
}

type BotService struct {
	chatProvider   ports.ChatAccessProvider
	docStore       ports.DocumentStoreProvider
	aiAgent        ports.AiAgentProvider
	projectService *ProjectService
	docService     *DocumentationService
	handlers       map[domain.MessageType]MessageHandler
}

func NewBotService(
	chat ports.ChatAccessProvider,
	docs ports.DocumentStoreProvider,
	ai ports.AiAgentProvider,
	ps *ProjectService,
	ds *DocumentationService,
) *BotService {
	if chat == nil {
		panic("chat provider cannot be nil")
	}
	if docs == nil {
		panic("document store cannot be nil")
	}
	if ai == nil {
		panic("AI agent cannot be nil")
	}
	if ps == nil {
		panic("project service cannot be nil")
	}
	if ds == nil {
		panic("documentation service cannot be nil")
	}

	base := baseHandler{
		docService:   ds,
		chatProvider: chat,
	}

	handlers := map[domain.MessageType]MessageHandler{
		domain.MessageTypeIdea:     &ideaHandler{base},
		domain.MessageTypeDecision: &decisionHandler{base},
		domain.MessageTypeStatus:   &statusHandler{base},
		domain.MessageTypeUnknown:  &unknownHandler{base},
	}

	return &BotService{
		chatProvider:   chat,
		docStore:       docs,
		aiAgent:        ai,
		projectService: ps,
		docService:     ds,
		handlers:       handlers,
	}
}

func (s *BotService) ProcessMessage(ctx context.Context, msg *domain.Message) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if msg == nil {
		return fmt.Errorf("message cannot be nil")
	}

	analysis, err := s.analyzeMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to analyze message: %w", err)
	}

	s.updateMessageWithAnalysis(msg, analysis)

	if !msg.HasReferences() {
		if err := s.detectAndAddReferences(ctx, msg); err != nil {
			return err
		}
	}

	handler := s.handlers[analysis.MessageType()]
	if suggestHandler, ok := handler.(SuggestionsHandler); ok {
		return suggestHandler.HandleWithAnalysis(ctx, msg, analysis)
	}
	return handler.Handle(ctx, msg)
}

func (s *BotService) analyzeMessage(ctx context.Context, msg *domain.Message) (*domain.MessageAnalysisResult, error) {
	result, err := s.aiAgent.AnalyzeMessage(ctx, msg.Content().Text())
	if err != nil {
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}
	return result, nil
}

func (s *BotService) updateMessageWithAnalysis(msg *domain.Message, analysis *domain.MessageAnalysisResult) {
	msg.UpdateCategory(analysis.Category())
	for _, ref := range analysis.References() {
		msg.AddReference(ref)
	}
}

func (s *BotService) detectAndAddReferences(ctx context.Context, msg *domain.Message) error {
	refs, err := s.aiAgent.DetectReferences(ctx, msg.Content().Text())
	if err != nil {
		return fmt.Errorf("failed to detect references: %w", err)
	}
	for _, ref := range refs {
		msg.AddReference(ref)
	}
	return nil
}
