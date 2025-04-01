package model

type RoutingPolicy struct {
	DefinedSets       DefinedSets       `json:"defined-sets"`
	PolicyDefinitions PolicyDefinitions `json:"policy-definitions"`
}

type DefinedSets struct {
	PrefixSets     PrefixSets     `json:"prefix-sets"`
	BGPDefinedSets BGPDefinedSets `json:"bgp-defined-sets"`
}

type PrefixSets struct {
	PrefixSet map[string]PrefixSet `json:"prefix-set"`
}

type PrefixSet struct {
	Name     string          `json:"name"`
	Config   PrefixSetConfig `json:"config"`
	Prefixes Prefixes        `json:"prefixes"`
}

type PrefixSetConfig struct {
	Name string `json:"name"`
	Mode string `json:"mode"`
}

type Prefixes struct {
	Prefix map[string]Prefix `json:"prefix"`
}

type Prefix struct {
	IPPrefix        string       `json:"ip-prefix"`
	MasklengthRange string       `json:"masklength-range"`
	Config          PrefixConfig `json:"config"`
}

type PrefixConfig struct {
	IPPrefix        string `json:"ip-prefix"`
	MasklengthRange string `json:"masklength-range"`
}

type BGPDefinedSets struct {
	CommunitySets    CommunitySets    `json:"community-sets"`
	ExtCommunitySets ExtCommunitySets `json:"ext-community-sets"`
}

type CommunitySets struct {
	CommunitySet map[string]CommunitySet `json:"community-set"`
}
type CommunitySet struct {
	CommunitySetName string             `json:"community-set-name"`
	Config           CommunitySetConfig `json:"config"`
}

type CommunitySetConfig struct {
	CommunitySetName string   `json:"community-set-name"`
	ReplaceConfig    bool     `json:"replace-config"`
	CommunityMember  []string `json:"community-member"`
}

type ExtCommunitySets struct {
	ExtCommunitySet map[string]ExtCommunitySet `json:"ext-community-set"`
}
type ExtCommunitySet struct {
	ExtCommunitySetName string                `json:"ext-community-set-name"`
	Config              ExtCommunitySetConfig `json:"config"`
}

type ExtCommunitySetConfig struct {
	ExtCommunitySetName string   `json:"ext-community-set-name"`
	ReplaceConfig       bool     `json:"replace-config"`
	ExtCommunityMember  []string `json:"ext-community-member"`
}

type PolicyDefinitions struct {
	PolicyDefinition map[string]PolicyDefinition `json:"policy-definition"`
}
type PolicyDefinition struct {
	Name       string                 `json:"name"`
	Config     PolicyDefinitionConfig `json:"config"`
	Statements Statements             `json:"statements"`
}

type PolicyDefinitionConfig struct {
	Name          string `json:"name"`
	DefaultAction string `json:"default-action"`
}

type Statements struct {
	Statement map[string]Statement `json:"statement"`
}
type Statement struct {
	Name       string          `json:"name"`
	Config     StatementConfig `json:"config"`
	Actions    Actions         `json:"actions"`
	Conditions Conditions      `json:"conditions"`
}

type StatementConfig struct {
	StatementName string `json:"statement-name"`
	Name          string `json:"name"`
	Seq           int    `json:"seq"`
	Description   string `json:"description"`
}

type Actions struct {
	Config     ActionsConfig `json:"config"`
	BGPActions BGPActions    `json:"bgp-actions"`
}

type ActionsConfig struct {
	PolicyResult string `json:"policy-result"`
}

type Conditions struct {
	BGPConditions BGPConditionsConfig `json:"bgp-conditions"`
}

type BGPConditionsConfig struct {
	Config ConditionsConfig `json:"config"`
}

type ConditionsConfig struct {
	ExtCommunitySet string   `json:"ext-community-set"`
	CallPolicies    []string `json:"call-policies"`
	CommunitySet    string   `json:"community-set"`
}

type BGPActions struct {
	Config          BGPActionsConfig `json:"config"`
	SetCommunity    SetCommunity     `json:"set-community"`
	SetExtCommunity SetExtCommunity  `json:"set-ext-community"`
}

type BGPActionsConfig struct {
	SetMed int `json:"set-med"`
}

type SetCommunity struct {
	Config    SetCommunityConfig `json:"config"`
	Reference Reference          `json:"reference"`
}

type Reference struct {
	Config ReferenceConfig `json:"config"`
}

type ReferenceConfig struct {
	CommunitySetRef string `json:"community-set-ref"`
}

type SetCommunityConfig struct {
	Method  string `json:"method"`
	Options string `json:"options"`
}

type SetExtCommunity struct {
	Config    SetExtCommunityConfig `json:"config"`
	Reference Reference             `json:"reference"`
}

type SetExtCommunityConfig struct {
	Method  string `json:"method"`
	Options string `json:"options"`
}
