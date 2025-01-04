package handlers

import (
	"context"
	"github.com/massimo-ua/quill/internal/domain"
	"github.com/massimo-ua/quill/internal/domain/services"
)

type MessageHandler struct {
	botService *services.BotService
}

func NewMessageHandler(bot *services.BotService) *MessageHandler {
	return &MessageHandler{botService: bot}
}

func (h *MessageHandler) HandleMessage(ctx context.Context, msg *domain.Message) error {
	return h.botService.ProcessMessage(ctx, msg)
}
