package main

type WTProject struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CurrentNode *int   `json:"currentNode,omitempty"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
}

type WTNode struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"projectId"`
	ParentID  *int   `json:"parentId,omitempty"`
	Name      string `json:"name"`
	Desc      string `json:"desc,omitempty"`
}

type WTBranchTag struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parentId"`
	ChildID  int    `json:"childId"`
	Name     string `json:"name"`
	Desc     string `json:"desc,omitempty"`
}

type WTFolder struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WTCard struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Desc          string `json:"desc,omitempty"`
	KeyWord       string `json:"keyWord"`
	ImageWord     string `json:"imageWord"`
	TriggerSystem bool   `json:"triggerSystem"`
	TriggerUser   bool   `json:"triggerUser"`
	TriggerAI     bool   `json:"triggerAi"`
}

type WTBlock struct {
	ID    int    `json:"id"`
	Title string `json:"title,omitempty"`
	Zone  string `json:"zone"`
}

type WTLine struct {
	ID        int    `json:"id"`
	SN        string `json:"sn"`
	ProjectID int    `json:"projectId"`
	Content   string `json:"content,omitempty"`
	Position  *int   `json:"position,omitempty"`
}

type NodeDetailBlock struct {
	Title string            `json:"title"`
	Lines map[string]string `json:"lines"`
}

type NodeDetailTrigger struct {
	Mode   string   `json:"mode"`
	Words  []string `json:"words"`
	System bool     `json:"system"`
	User   bool     `json:"user"`
	AI     bool     `json:"ai"`
}

type NodeDetailCard struct {
	Name    string                     `json:"name"`
	Desc    string                     `json:"desc,omitempty"`
	Trigger *NodeDetailTrigger         `json:"trigger,omitempty"`
	Content map[string]NodeDetailBlock `json:"content"`
	Image   []string                   `json:"image,omitempty"`
}

type NodeDetailFolder struct {
	Name    string                      `json:"name"`
	Cards   map[string]NodeDetailCard   `json:"cards,omitempty"`
	Folders map[string]NodeDetailFolder `json:"folders,omitempty"`
}

type NodeDetail struct {
	NodeID int                         `json:"nodeId"`
	Name   string                      `json:"name"`
	Desc   string                      `json:"desc,omitempty"`
	Struct map[string]NodeDetailFolder `json:"struct"`
}

type ImmediateChange struct {
	Name   string `json:"name"`
	Action string `json:"action"`
	Level  string `json:"level"`
	Target *int   `json:"target"`
}

type SaveFolderChange struct {
	Name *string `json:"name,omitempty"`
	Desc *string `json:"desc,omitempty"`
}

type SaveCardChange struct {
	Name    *string            `json:"name,omitempty"`
	Desc    *string            `json:"desc,omitempty"`
	Trigger *NodeDetailTrigger `json:"trigger,omitempty"`
	Image   []string           `json:"image,omitempty"`
}

type SaveBlockChange struct {
	Name    *string `json:"name,omitempty"`
	Content *string `json:"content,omitempty"`
}

type SaveChanges struct {
	Folder map[string]SaveFolderChange `json:"folder,omitempty"`
	Card   map[string]SaveCardChange   `json:"card,omitempty"`
	Block  map[string]SaveBlockChange  `json:"block,omitempty"`
}

type WTNodeChange struct {
	ID           int  `json:"id"`
	Action       string `json:"action"`
	Level        string `json:"level"`
	NodeID       int  `json:"nodeId"`
	Target       *int `json:"target,omitempty"`
	DetailFolder *int `json:"detailFolder,omitempty"`
	DetailCard   *int `json:"detailCard,omitempty"`
	DetailBlock  *int `json:"detailBlock,omitempty"`
	DetailLine   *int `json:"detailLine,omitempty"`
}

type WTApp struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WTConversation struct {
	ID          int    `json:"id"`
	AppID       int    `json:"appId"`
	Name        string `json:"name"`
	CurrentNode *int   `json:"currentNode,omitempty"`
}

type WTDialogue struct {
	ID                   int    `json:"id"`
	ConversationID       int    `json:"conversationId"`
	CreateTime           string `json:"createTime"`
	RequestContent       string `json:"requestContent"`
	ResponseContent      string `json:"responseContent"`
	RequestSystemPrompt  string `json:"requestSystemPrompt,omitempty"`
	ResponseSystemPrompt string `json:"responseSystemPrompt,omitempty"`
	NodeID               *int   `json:"nodeId,omitempty"`
	RequestPoint         *int   `json:"requestPoint,omitempty"`
	ResponsePoint        *int   `json:"responsePoint,omitempty"`
	RequestToken         *int   `json:"requestToken,omitempty"`
	ResponseToken        *int   `json:"responseToken,omitempty"`
}

type WTDialogueImage struct {
	ID         int    `json:"id"`
	DialogueID int    `json:"dialogueId"`
	ImageURL   string `json:"imageUrl,omitempty"`
	ImagePath  string `json:"imagePath,omitempty"`
	Prompt     string `json:"prompt"`
}

type ProjectInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	FileName  string `json:"fileName"`
	Type      string `json:"type"`
	UpdatedAt string `json:"updatedAt"`
}

type Project struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	FileName    string         `json:"fileName"`
	Type        string         `json:"type"`
	CurrentNode *int           `json:"currentNode,omitempty"`
	CreateTime  string         `json:"createTime"`
	UpdateTime  string         `json:"updateTime"`
	Timeline    []NodeInfo     `json:"timeline"`
}

type NodeInfo struct {
	ID        int             `json:"id"`
	ParentID  *int            `json:"parentId,omitempty"`
	Name      string          `json:"name"`
	Note      string          `json:"note,omitempty"`
	Tags      []string        `json:"tags"`
	CreatedAt string          `json:"createdAt,omitempty"`
	PreText   []PromptEntry   `json:"pre_text"`
	PostText  []PromptEntry   `json:"post_text"`
	PrePrompt []PromptEntry   `json:"pre_prompt"`
	WorldBook []WorldBookEntry `json:"world_book"`
}

type ReplayedFolder struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	ChangeID int             `json:"changeId"`
	Cards    []ReplayedCard  `json:"cards"`
}

type ReplayedCard struct {
	ID            int              `json:"id"`
	Name          string           `json:"name"`
	Desc          string           `json:"desc,omitempty"`
	KeyWord       string           `json:"keyWord"`
	TriggerSystem bool             `json:"triggerSystem"`
	TriggerUser   bool             `json:"triggerUser"`
	TriggerAI     bool             `json:"triggerAi"`
	ChangeID      int              `json:"changeId"`
	Blocks        []ReplayedBlock  `json:"blocks"`
}

type ReplayedBlock struct {
	ID       int            `json:"id"`
	Title    string         `json:"title,omitempty"`
	Zone     string         `json:"zone"`
	ChangeID int            `json:"changeId"`
	Lines    []ReplayedLine `json:"lines"`
}

type ReplayedLine struct {
	ID       int     `json:"id"`
	SN       string  `json:"sn"`
	Content  *string `json:"content"`
	SyncDot  string  `json:"syncDot"`
	ChangeID int     `json:"changeId"`
}

type NodeContent struct {
	NodeID  int              `json:"nodeId"`
	Folders []ReplayedFolder `json:"folders"`
	Pre     []ReplayedBlock  `json:"pre"`
	Post    []ReplayedBlock  `json:"post"`
	Global  []ReplayedBlock  `json:"global"`
}

type LineMarker struct {
	Serial       string `json:"serial"`
	Type         string `json:"type"`
	Content      string `json:"content,omitempty"`
	InsertBefore string `json:"insertBefore,omitempty"`
}

type ContentItem struct {
	ID          string `json:"id,omitempty"`
	Serial      string `json:"serial,omitempty"`
	Title       string `json:"title,omitempty"`
	Content     string `json:"content,omitempty"`
	KeySystem   bool   `json:"keySystem,omitempty"`
	KeyUser     bool   `json:"keyUser,omitempty"`
	KeyAI       bool   `json:"keyAI,omitempty"`
	KeyRegion   int    `json:"keyRegion,omitempty"`
	ValueRegion int    `json:"valueRegion,omitempty"`
	Collapsed   bool   `json:"collapsed,omitempty"`
}

type WorldBookEntry struct {
	ID             string        `json:"id"`
	Key            string        `json:"key"`
	Value          string        `json:"value"`
	ParentID       string        `json:"parentId,omitempty"`
	IsFolder       bool          `json:"isFolder,omitempty"`
	Name           string        `json:"name,omitempty"`
	MatchMode      string        `json:"matchMode,omitempty"`
	KeyRegion      int           `json:"keyRegion,omitempty"`
	ValueRegion    int           `json:"valueRegion,omitempty"`
	Group          string        `json:"group,omitempty"`
	ContentItems   []ContentItem `json:"contentItems,omitempty"`
	Keywords       []string      `json:"keywords,omitempty"`
	T2iKeywords    []string      `json:"t2iKeywords,omitempty"`
	SourceEntryID  string        `json:"sourceEntryId,omitempty"`
	LineMarkers    []LineMarker  `json:"lineMarkers,omitempty"`
	DeletedItemIDs []string      `json:"deletedItemIds,omitempty"`
	LocalItems     []ContentItem `json:"localItems,omitempty"`
	ItemOverrides  map[string]map[string]any `json:"itemOverrides,omitempty"`
}

type SyncRule struct {
	TargetNodeID  string `json:"targetNodeId"`
	TargetEntryID string `json:"targetEntryId"`
	Mode          string `json:"mode"`
}

type PromptEntry struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Content       string `json:"content"`
	Enabled       bool   `json:"enabled"`
	SourceEntryID string `json:"sourceEntryId,omitempty"`
}

type NodeBinding struct {
	AppID          string `json:"appId"`
	AppName        string `json:"appName"`
	ConversationID string `json:"conversationId"`
	ConfigID       string `json:"configId"`
}

type TimelineNode struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Note      string           `json:"note,omitempty"`
	Tags      []string         `json:"tags"`
	ParentID  string           `json:"parentId,omitempty"`
	CreatedAt string           `json:"createdAt"`
	PreText   []PromptEntry    `json:"pre_text"`
	PostText  []PromptEntry    `json:"post_text"`
	PrePrompt []PromptEntry    `json:"pre_prompt"`
	WorldBook []WorldBookEntry `json:"world_book"`
	Binding   *NodeBinding     `json:"binding,omitempty"`
}

type Config struct {
	Theme              string `json:"theme"`
	Language           string `json:"language"`
	LastOpenedProject  string `json:"lastOpenedProject"`
	ColorScheme        int    `json:"colorScheme"`
	ColorMode          string `json:"colorMode"`
	SystemPrompt       string `json:"systemPrompt"`
	SystemPromptType   string `json:"systemPromptType"`
	DebugMode          bool   `json:"debugMode"`
	SafeMode           bool   `json:"safeMode"`
	SafeModeAction     string `json:"safeModeAction"`
	SafeModeTemplate   string `json:"safeModeTemplate"`
	DebugTestReply     string `json:"debugTestReply"`
	BytedanceApiKey    string `json:"bytedanceApiKey"`
	AutoGenerateImage  bool   `json:"autoGenerateImage"`
	NoImageMode        bool   `json:"noImageMode"`
	StrictMode         bool   `json:"strictMode"`
	ClaudeApiKey       string `json:"claudeApiKey"`
	GeminiApiKey       string `json:"geminiApiKey"`
	GrokApiKey         string `json:"grokApiKey"`
}

type SessionState struct {
	OpenTabs           []OpenTabState     `json:"openTabs"`
	ActiveTabIndex     int                `json:"activeTabIndex"`
	ActiveConversation *ConversationState `json:"activeConversation"`
}

type OpenTabState struct {
	ProjectFileName string `json:"projectFileName"`
	CurrentNodeId   string `json:"currentNodeId"`
}

type ConversationState struct {
	AppId          string `json:"appId"`
	ConversationId string `json:"conversationId"`
}

type ByteDanceImageData struct {
	URL string `json:"url"`
}

type ByteDanceImageError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ByteDanceImageResponse struct {
	Created int                    `json:"created"`
	Data    []ByteDanceImageData   `json:"data"`
	Error   *ByteDanceImageError   `json:"error"`
}

