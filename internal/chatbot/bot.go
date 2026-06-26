package chatbot

import "strings"

type Bot struct{}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) Respond(message string) string {
	msg := strings.ToLower(message)

	switch {
	case strings.Contains(msg, "senha"):
		return "Para redefinir sua senha, acesse a página de login e clique em 'Esqueci minha senha'."
	case strings.Contains(msg, "cobrança") || strings.Contains(msg, "fatura") || strings.Contains(msg, "pagamento"):
		return "Questões de cobrança são tratadas pelo nosso time financeiro. Um atendente humano vai analisar seu chamado em breve."
	case strings.Contains(msg, "horário") || strings.Contains(msg, "atendimento"):
		return "Nosso atendimento funciona de segunda a sexta, das 9h às 18h."
	case strings.Contains(msg, "cancelar"):
		return "Sentimos muito que você queira cancelar. Um de nossos atendentes vai entrar em contato para entender melhor o motivo."
	case strings.Contains(msg, "obrigado") || strings.Contains(msg, "obrigada"):
		return "Por nada! Estamos aqui para ajudar sempre que precisar."
	case strings.Contains(msg, "oi") || strings.Contains(msg, "olá") || strings.Contains(msg, "ola"):
		return "Olá! Sou o assistente virtual de suporte. Como posso te ajudar hoje?"
	default:
		return "Recebemos sua mensagem e em breve um atendente vai te responder. Enquanto isso, você pode nos contar mais detalhes sobre o problema."
	}
}