package cli

import (
	"context"
	"errors"
	"io"
	"strings"
	"unicode"
)

type CommandActions struct {
	Clear func(); Status func() StartupStatus; Models func(context.Context) ([]string,error); Model func(context.Context,string) error; Studio func(context.Context,string) error; Tools func(context.Context)([]string,error); History func(context.Context)([]string,error); Diff func(context.Context)(string,error); Undo func(context.Context) error; Config func(context.Context)(string,error)
}

func NewCommandHandler(out io.Writer, next InputHandler) InputHandler { return NewCommandHandlerWithActions(out, CommandActions{}, next) }
func NewCommandHandlerWithClear(out io.Writer, clear func(), next InputHandler) InputHandler { return NewCommandHandlerWithActions(out, CommandActions{Clear:clear}, next) }
func safeModelIDs(models []string) []string { safe:=make([]string,0,len(models)); for _,model:=range models { unsafe:=false; for _,r:=range model { if unicode.IsControl(r)||unicode.Is(unicode.Cf,r){unsafe=true;break} }; if !unsafe {safe=append(safe,model)} }; return safe }
func safeDisplayText(text string) string { return strings.Map(func(r rune) rune { if unicode.IsControl(r)||unicode.Is(unicode.Cf,r){return -1}; return r },text) }
func safeMultilineDisplayText(text string) string { return strings.Map(func(r rune) rune { if r=='\n'||r=='\t' { return r }; if unicode.IsControl(r)||unicode.Is(unicode.Cf,r){return -1}; return r },text) }

// WriteSafeMultiline writes terminal-facing text while preserving intentional
// newlines and tabs and removing control/format characters that could alter
// terminal state or hide content.
func WriteSafeMultiline(out io.Writer, text string) error {
	if out == nil { return errors.New("cli: nil output") }
	safe := safeMultilineDisplayText(text)
	n, err := io.WriteString(out, safe)
	if err == nil && n != len(safe) { return io.ErrShortWrite }
	return err
}

func NewCommandHandlerWithActions(out io.Writer, actions CommandActions, next InputHandler) InputHandler {
	arenaConnected:=false
	return func(ctx context.Context,input Input)(bool,error){
		if input.Kind==InputCommand&&input.Command=="help" { if out==nil{return false,errors.New("cli: nil output")}; _,err:=io.WriteString(out,HelpText()); return false,err }
		if input.Kind==InputCommand&&input.Command=="exit" {return true,nil}
		if input.Kind==InputCommand&&input.Command=="clear"&&actions.Clear!=nil {actions.Clear();return false,nil}
		if input.Kind==InputCommand&&input.Command=="status"&&actions.Status!=nil {if out==nil{return false,errors.New("cli: nil output")}; status:=actions.Status();if arenaConnected{status.Arena="connected"};_,err:=io.WriteString(out,StatusText(status));return false,err}
		if input.Kind==InputCommand&&input.Command=="models"&&actions.Models!=nil {models,err:=actions.Models(ctx);if err!=nil{return false,err};models=safeModelIDs(models);arenaConnected=true;if out==nil{return false,errors.New("cli: nil output")};if len(models)==0{_,err=io.WriteString(out,"No Arena models available.\n");return false,err};_,err=io.WriteString(out,strings.Join(models,"\n")+"\n");return false,err}
		if input.Kind==InputCommand&&input.Command=="model"&&strings.TrimSpace(input.Argument)==""&&actions.Models!=nil {models,err:=actions.Models(ctx);if err!=nil{if actions.Model!=nil{return false,actions.Model(ctx,input.Argument)};return false,err};models=safeModelIDs(models);arenaConnected=true;if out==nil{return false,errors.New("cli: nil output")};if len(models)==0{_,err=io.WriteString(out,"No Arena models available.\n");return false,err};if _,err=io.WriteString(out,"Select a model with /model <id>:\n");err!=nil{return false,err};_,err=io.WriteString(out,strings.Join(models,"\n")+"\n");return false,err}
		if input.Kind==InputCommand&&input.Command=="model"&&actions.Model!=nil {err:=actions.Model(ctx,input.Argument);if err==nil{arenaConnected=true};return false,err}
		if input.Kind==InputCommand&&input.Command=="studio"&&actions.Studio!=nil{return false,actions.Studio(ctx,input.Argument)}
		if input.Kind==InputCommand&&input.Command=="tools"&&actions.Tools!=nil {tools,err:=actions.Tools(ctx);if err!=nil{return false,err};if out==nil{return false,errors.New("cli: nil output")};if len(tools)==0{_,err=io.WriteString(out,"No MCP tools available.\n");return false,err};safeTools:=make([]string,len(tools));for i:=range tools{safeTools[i]=safeDisplayText(tools[i])};_,err=io.WriteString(out,strings.Join(safeTools,"\n")+"\n");return false,err}
		if input.Kind==InputCommand&&input.Command=="history"&&actions.History!=nil {history,err:=actions.History(ctx);if err!=nil{return false,err};if out==nil{return false,errors.New("cli: nil output")};if len(history)==0{_,err=io.WriteString(out,"No session history recorded.\n");return false,err};safeHistory:=make([]string,len(history));for i:=range history{safeHistory[i]=safeDisplayText(history[i])};_,err=io.WriteString(out,strings.Join(safeHistory,"\n")+"\n");return false,err}
		if input.Kind==InputCommand&&input.Command=="diff"&&actions.Diff!=nil {diff,err:=actions.Diff(ctx);if err!=nil{return false,err};if diff==""{return false,nil};return false,WriteSafeMultiline(out,diff)}
		if input.Kind==InputCommand&&input.Command=="undo" {if actions.Undo!=nil{return false,actions.Undo(ctx)};if out==nil{return false,errors.New("cli: nil output")};_,err:=io.WriteString(out,"Nothing to undo.\n");return false,err}
		if input.Kind==InputCommand&&input.Command=="config"&&actions.Config!=nil {config,err:=actions.Config(ctx);if err!=nil{return false,err};if config==""{return false,nil};return false,WriteSafeMultiline(out,config)}
		if next==nil{return false,nil};return next(ctx,input)
	}
}
