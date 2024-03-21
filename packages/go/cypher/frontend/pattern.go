// Copyright 2023 Specter Ops, Inc.
//
// Licensed under the Apache License, Version 2.0
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package frontend

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"github.com/specterops/bloodhound/cypher/model/cypher"
	"github.com/specterops/bloodhound/cypher/parser"
	"github.com/specterops/bloodhound/dawgs/graph"
	"strconv"
)

type WhereVisitor struct {
	BaseVisitor

	Where *cypher.Where
}

func NewWhereVisitor() *WhereVisitor {
	return &WhereVisitor{
		Where: cypher.NewWhere(),
	}
}

func (s *WhereVisitor) EnterOC_Expression(ctx *parser.OC_ExpressionContext) {
	s.ctx.Enter(NewExpressionVisitor())
}

func (s *WhereVisitor) ExitOC_Expression(ctx *parser.OC_ExpressionContext) {
	s.Where.Add(s.ctx.Exit().(*ExpressionVisitor).Expression)
}

type NodePatternVisitor struct {
	BaseVisitor

	NodePattern *cypher.NodePattern
}

func (s *NodePatternVisitor) EnterOC_Variable(ctx *parser.OC_VariableContext) {
	s.ctx.Enter(NewVariableVisitor())
}

func (s *NodePatternVisitor) ExitOC_Variable(ctx *parser.OC_VariableContext) {
	s.NodePattern.Binding = s.ctx.Exit().(*VariableVisitor).Variable
}

func (s *NodePatternVisitor) EnterOC_NodeLabels(ctx *parser.OC_NodeLabelsContext) {
}

func (s *NodePatternVisitor) ExitOC_NodeLabels(ctx *parser.OC_NodeLabelsContext) {
}

func (s *NodePatternVisitor) EnterOC_NodeLabel(ctx *parser.OC_NodeLabelContext) {
}

func (s *NodePatternVisitor) ExitOC_NodeLabel(ctx *parser.OC_NodeLabelContext) {
}

func (s *NodePatternVisitor) EnterOC_LabelName(ctx *parser.OC_LabelNameContext) {
	s.ctx.Enter(&SymbolicNameOrReservedWordVisitor{})
}

func (s *NodePatternVisitor) ExitOC_LabelName(ctx *parser.OC_LabelNameContext) {
	kind := graph.StringKind(s.ctx.Exit().(*SymbolicNameOrReservedWordVisitor).Name)
	s.NodePattern.Kinds = append(s.NodePattern.Kinds, kind)
}

func (s *NodePatternVisitor) EnterOC_Properties(ctx *parser.OC_PropertiesContext) {
	s.ctx.Enter(NewPropertiesVisitor())
}

func (s *NodePatternVisitor) ExitOC_Properties(ctx *parser.OC_PropertiesContext) {
	s.NodePattern.Properties = s.ctx.Exit().(*PropertiesVisitor).Properties
}

type RelationshipPatternVisitor struct {
	BaseVisitor

	RelationshipPattern *cypher.RelationshipPattern
}

func (s *RelationshipPatternVisitor) EnterOC_RelationshipTypes(ctx *parser.OC_RelationshipTypesContext) {
}

func (s *RelationshipPatternVisitor) ExitOC_RelationshipTypes(ctx *parser.OC_RelationshipTypesContext) {
}

func (s *RelationshipPatternVisitor) EnterOC_Dash(ctx *parser.OC_DashContext) {
}

func (s *RelationshipPatternVisitor) ExitOC_Dash(ctx *parser.OC_DashContext) {
}

func (s *RelationshipPatternVisitor) EnterOC_RelationshipDetail(ctx *parser.OC_RelationshipDetailContext) {
}

func (s *RelationshipPatternVisitor) ExitOC_RelationshipDetail(ctx *parser.OC_RelationshipDetailContext) {
}

func (s *RelationshipPatternVisitor) EnterOC_RelTypeName(ctx *parser.OC_RelTypeNameContext) {
	s.ctx.Enter(&SymbolicNameOrReservedWordVisitor{})
}

func (s *RelationshipPatternVisitor) ExitOC_RelTypeName(ctx *parser.OC_RelTypeNameContext) {
	kind := graph.StringKind(s.ctx.Exit().(*SymbolicNameOrReservedWordVisitor).Name)
	s.RelationshipPattern.Kinds = append(s.RelationshipPattern.Kinds, kind)
}

func (s *RelationshipPatternVisitor) EnterOC_Variable(ctx *parser.OC_VariableContext) {
	s.ctx.Enter(NewVariableVisitor())
}

func (s *RelationshipPatternVisitor) ExitOC_Variable(ctx *parser.OC_VariableContext) {
	s.RelationshipPattern.Binding = s.ctx.Exit().(*VariableVisitor).Variable
}

func (s *RelationshipPatternVisitor) EnterOC_LeftArrowHead(ctx *parser.OC_LeftArrowHeadContext) {
	s.RelationshipPattern.Direction = graph.DirectionInbound
}

func (s *RelationshipPatternVisitor) ExitOC_LeftArrowHead(ctx *parser.OC_LeftArrowHeadContext) {
}

func (s *RelationshipPatternVisitor) EnterOC_RightArrowHead(ctx *parser.OC_RightArrowHeadContext) {
	if s.RelationshipPattern.Direction == graph.DirectionInbound {
		s.RelationshipPattern.Direction = graph.DirectionBoth
	} else {
		s.RelationshipPattern.Direction = graph.DirectionOutbound
	}
}

func (s *RelationshipPatternVisitor) ExitOC_RightArrowHead(ctx *parser.OC_RightArrowHeadContext) {
}

func (s *RelationshipPatternVisitor) EnterOC_RangeLiteral(ctx *parser.OC_RangeLiteralContext) {
	const (
		stateStart int = iota
		stateFirstIndex
		stateSecondIndex
	)

	// Create a new relationship pattern range for the relationship pattern being built
	s.RelationshipPattern.Range = &cypher.PatternRange{}

	// Start at the start state for the mini-parser below
	state := stateStart

	for _, tokenLeaf := range ctx.GetChildren() {
		switch typedTokenLeaf := tokenLeaf.(type) {
		case *antlr.TerminalNodeImpl:
			switch typedTokenLeaf.GetSymbol().GetTokenType() {
			case TokenTypeAsterisk:
				state = stateFirstIndex

			case TokenTypeRange:
				state = stateSecondIndex

			default:
				s.ctx.AddErrors(fmt.Errorf("unexpected token in pattern range: %s", typedTokenLeaf.GetText()))
			}

		case *parser.OC_IntegerLiteralContext:
			if value, err := strconv.ParseInt(typedTokenLeaf.GetText(), 10, 64); err != nil {
				s.ctx.AddErrors(fmt.Errorf("failed parsing range literal: %w", err))
			} else {
				switch state {
				case stateFirstIndex:
					s.RelationshipPattern.Range.StartIndex = &value

				case stateSecondIndex:
					s.RelationshipPattern.Range.EndIndex = &value

				default:
					s.ctx.AddErrors(fmt.Errorf("invalid integer literal state: %d", state))
				}
			}
		}
	}
}

func (s *RelationshipPatternVisitor) ExitOC_RangeLiteral(ctx *parser.OC_RangeLiteralContext) {
}

func (s *RelationshipPatternVisitor) EnterOC_Properties(ctx *parser.OC_PropertiesContext) {
	s.ctx.Enter(NewPropertiesVisitor())
}

func (s *RelationshipPatternVisitor) ExitOC_Properties(ctx *parser.OC_PropertiesContext) {
	s.RelationshipPattern.Properties = s.ctx.Exit().(*PropertiesVisitor).Properties
}

type PatternPredicateVisitor struct {
	BaseVisitor

	PatternPredicate *cypher.PatternPredicate
}

func NewPatternPredicateVisitor() *PatternPredicateVisitor {
	return &PatternPredicateVisitor{
		PatternPredicate: cypher.NewPatternPredicate(),
	}
}

func (s *PatternPredicateVisitor) EnterOC_NodePattern(ctx *parser.OC_NodePatternContext) {
	s.ctx.Enter(&NodePatternVisitor{
		NodePattern: &cypher.NodePattern{},
	})
}

func (s *PatternPredicateVisitor) ExitOC_NodePattern(ctx *parser.OC_NodePatternContext) {
	s.PatternPredicate.AddElement(s.ctx.Exit().(*NodePatternVisitor).NodePattern)
}

func (s *PatternPredicateVisitor) EnterOC_RelationshipPattern(ctx *parser.OC_RelationshipPatternContext) {
	s.ctx.Enter(&RelationshipPatternVisitor{
		RelationshipPattern: &cypher.RelationshipPattern{
			Direction: graph.DirectionBoth,
		},
	})
}

func (s *PatternPredicateVisitor) ExitOC_RelationshipPattern(ctx *parser.OC_RelationshipPatternContext) {
	s.PatternPredicate.AddElement(s.ctx.Exit().(*RelationshipPatternVisitor).RelationshipPattern)
}

type PatternVisitor struct {
	BaseVisitor

	currentPart  *cypher.PatternPart
	PatternParts []*cypher.PatternPart
}

func (s *PatternVisitor) EnterOC_AnonymousPatternPart(ctx *parser.OC_AnonymousPatternPartContext) {
}

func (s *PatternVisitor) ExitOC_AnonymousPatternPart(ctx *parser.OC_AnonymousPatternPartContext) {
}

func (s *PatternVisitor) EnterOC_PatternElementChain(ctx *parser.OC_PatternElementChainContext) {
}

func (s *PatternVisitor) ExitOC_PatternElementChain(ctx *parser.OC_PatternElementChainContext) {
}

func (s *PatternVisitor) EnterOC_PatternElement(ctx *parser.OC_PatternElementContext) {
}

func (s *PatternVisitor) ExitOC_PatternElement(ctx *parser.OC_PatternElementContext) {
}

func (s *PatternVisitor) EnterOC_RelationshipsPattern(ctx *parser.OC_RelationshipsPatternContext) {
	s.currentPart = &cypher.PatternPart{}
}

func (s *PatternVisitor) ExitOC_RelationshipsPattern(ctx *parser.OC_RelationshipsPatternContext) {
	s.PatternParts = append(s.PatternParts, s.currentPart)
}

func (s *PatternVisitor) EnterOC_PatternPart(ctx *parser.OC_PatternPartContext) {
	s.currentPart = &cypher.PatternPart{}
}

func (s *PatternVisitor) ExitOC_PatternPart(ctx *parser.OC_PatternPartContext) {
	s.PatternParts = append(s.PatternParts, s.currentPart)
}

func (s *PatternVisitor) EnterOC_ShortestPathPattern(ctx *parser.OC_ShortestPathPatternContext) {
	if HasTokens(ctx, parser.CypherLexerSHORTESTPATH) {
		s.currentPart.ShortestPathPattern = true
	} else if HasTokens(ctx, parser.CypherLexerALLSHORTESTPATHS) {
		s.currentPart.AllShortestPathsPattern = true
	}
}

func (s *PatternVisitor) ExitOC_ShortestPathPattern(ctx *parser.OC_ShortestPathPatternContext) {
}

func (s *PatternVisitor) EnterOC_Variable(ctx *parser.OC_VariableContext) {
	s.ctx.Enter(NewVariableVisitor())
}

func (s *PatternVisitor) ExitOC_Variable(ctx *parser.OC_VariableContext) {
	s.currentPart.Binding = s.ctx.Exit().(*VariableVisitor).Variable
}

func (s *PatternVisitor) EnterOC_NodePattern(ctx *parser.OC_NodePatternContext) {
	s.ctx.Enter(&NodePatternVisitor{
		NodePattern: &cypher.NodePattern{},
	})
}

func (s *PatternVisitor) ExitOC_NodePattern(ctx *parser.OC_NodePatternContext) {
	s.currentPart.AddPatternElements(s.ctx.Exit().(*NodePatternVisitor).NodePattern)
}

func (s *PatternVisitor) EnterOC_RelationshipPattern(ctx *parser.OC_RelationshipPatternContext) {
	s.ctx.Enter(&RelationshipPatternVisitor{
		RelationshipPattern: &cypher.RelationshipPattern{
			Direction: graph.DirectionBoth,
		},
	})
}

func (s *PatternVisitor) ExitOC_RelationshipPattern(ctx *parser.OC_RelationshipPatternContext) {
	s.currentPart.AddPatternElements(s.ctx.Exit().(*RelationshipPatternVisitor).RelationshipPattern)
}
