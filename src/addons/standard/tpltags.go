package standard

import (
	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"io/ioutil"
	apptpl "tpl"
)

func (Extension) initTplTags() {
	pongo2.ReplaceTag("ssi", tagSSI)
}

type tagSSINode struct {
	filename string
	content  string
	template *pongo2.Template
}

func (node *tagSSINode) Execute(ctx *pongo2.ExecutionContext, writer pongo2.TemplateWriter) *pongo2.Error {
	if node.template != nil {
		// Execute the template within the current context
		includeCtx := make(pongo2.Context)
		includeCtx.Update(ctx.Public)
		includeCtx.Update(ctx.Private)
		content, err := node.template.Execute(includeCtx)
		if err != nil {
			return err.(*pongo2.Error)
		}
		writer.WriteString(content)
	} else {
		// Just print out the content
		writer.WriteString(node.content)
	}
	return nil
}

func tagSSI(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	SSINode := &tagSSINode{}

	if fileToken := arguments.MatchType(pongo2.TokenString); fileToken != nil {
		SSINode.filename = fileToken.Val

		if arguments.Match(pongo2.TokenIdentifier, "parsed") != nil {
			// parsed
			temporaryTpl, err := pongo2.DefaultSet.FromFile(fileToken.Val)
			if err != nil {
				return nil, err.(*pongo2.Error)
			}
			SSINode.template = temporaryTpl
		} else {
			// plaintext

			fileReader, err := apptpl.TplDefaultLoader.Get(fileToken.Val)

			if err != nil {

				logrus.WithField("_service", addonName).WithError(err).Warningf("pongo2: tag `ssi`, load file by name %q", fileToken.Val)
				SSINode.content = ""
				return SSINode, nil
				// return nil, (&pongo2.Error{
				// 	Sender:   "tag:ssi",
				// 	ErrorMsg: err.Error(),
				// })
			}

			buf, err := ioutil.ReadAll(fileReader)

			if err != nil {
				return nil, (&pongo2.Error{
					Sender:   "tag:ssi",
					ErrorMsg: err.Error(),
				})
			}

			SSINode.content = string(buf)
		}
	} else {
		return nil, arguments.Error("First argument must be a string.", nil)
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed SSI-tag argument.", nil)
	}

	return SSINode, nil
}
