package lib

import "github.com/bwmarrin/discordgo"

func GetModalDataValue(data *discordgo.ModalSubmitInteractionData, componentIndex int, inputIndex int) string {
	return data.Components[componentIndex].(*discordgo.ActionsRow).Components[inputIndex].(*discordgo.TextInput).Value
}

func GetOptionByName(options []*discordgo.ApplicationCommandInteractionDataOption, name string) *discordgo.ApplicationCommandInteractionDataOption {
	for _, option := range options {
		if option.Name == name {
			return option
		}
	}
	return nil
}

func SendEmptyInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})
}
