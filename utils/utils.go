package utils

import "github.com/anthropics/anthropic-sdk-go"

func MakeTextContent(text string) []anthropic.ContentBlockParamUnion {
	return []anthropic.ContentBlockParamUnion{{
		OfRequestTextBlock: &anthropic.TextBlockParam{
			Text: text,
		},
	}}
}
