package parser

import "time"

type FeedEntry struct {
	FeedEntryType       int                 `json:"feed_entry_type"`
	Rank                *int                `json:"rank"`
	UpdatedAt           time.Time           `json:"updated_at"`
	DestinationURL      string              `json:"destination_url"`
	Item                Item                `json:"item"`
	Modules             []Module            `json:"modules"`
	Children            []Child             `json:"children"`
	AnalyticsProperties AnalyticsProperties `json:"analytics_properties"`
	EntityID            int64               `json:"entity_id"`
	Destination         Destination         `json:"destination"`
	Category            string              `json:"category"`
	Page                string              `json:"page"`
	EventsToTrack       []string            `json:"events_to_track"`
	TrackableID         string              `json:"trackable_id"`
}

type Item struct {
	EntityType                  string        `json:"entity_type"`
	ID                          int64         `json:"id"`
	Name                        string        `json:"name"`
	Type                        string        `json:"type"`
	DeviceType                  int           `json:"device_type"`
	SponsoredActivity           bool          `json:"sponsored_activity"`
	PartnerTreatment            string        `json:"partner_treatment"`
	NoTreatmentBcOptOut         bool          `json:"no_treatment_bc_opt_out"`
	BasicIntegration            bool          `json:"basic_integration"`
	Private                     bool          `json:"private"`
	CommentCount                int           `json:"comment_count"`
	KudosCount                  int           `json:"kudos_count"`
	HasKudoed                   bool          `json:"has_kudoed"`
	AchievementCount            int           `json:"achievement_count"`
	Commute                     bool          `json:"commute"`
	BoostInFeed                 bool          `json:"boost_in_feed"`
	SportType                   string        `json:"sport_type"`
	Description                 string        `json:"description"`
	DescriptionMentionsMetadata []interface{} `json:"description_mentions_metadata"`
	Athlete                     Athlete       `json:"athlete"`
	Distance                    float64       `json:"distance"`
	Trainer                     bool          `json:"trainer"`
	MovingTime                  int           `json:"moving_time"`
	ElapsedTime                 int           `json:"elapsed_time"`
	StartDate                   time.Time     `json:"start_date"`
	BoundingBox                 [][]float64   `json:"bounding_box"`
}

type Athlete struct {
	ID int `json:"id"`
}

type Module struct {
	Type         string        `json:"type"`
	ModuleFields []ModuleField `json:"module_fields"`
	Submodules   []Submodule   `json:"submodules,omitempty"`
}

type ModuleField struct {
	Key                    string               `json:"key"`
	Value                  interface{}          `json:"value,omitempty"`
	ValueObject            interface{}          `json:"value_object,omitempty"`
	ItemKey                string               `json:"item_key,omitempty"`
	Destination            *Destination         `json:"destination,omitempty"`
	Element                string               `json:"element,omitempty"`
	ShouldTrackImpressions bool                 `json:"should_track_impressions,omitempty"`
	TrackableID            string               `json:"trackable_id,omitempty"`
	AnalyticsProperties    *AnalyticsProperties `json:"analytics_properties,omitempty"`
}

type ColorToken struct {
	Name         string       `json:"name"`
	CurrentValue CurrentValue `json:"current_value"`
}

type CurrentValue struct {
	LightHex string `json:"light_hex"`
	DarkHex  string `json:"dark_hex"`
}

type Destination struct {
	URL string `json:"url"`
}

type AnalyticsProperties struct {
	HasPartnerTreatment      bool          `json:"has_partner_treatment"`
	PartnerTreatmentType     string        `json:"partner_treatment_type"`
	OptedOut                 bool          `json:"opted_out"`
	MapTreatmentEligible     bool          `json:"map_treatment_eligible"`
	ImageTreatmentEligible   bool          `json:"image_treatment_eligible"`
	BannerCta                bool          `json:"banner_cta"`
	InGroup                  bool          `json:"in_group"`
	Owner                    bool          `json:"owner"`
	AchievementCount         int           `json:"achievement_count"`
	MentionedAthletes        []interface{} `json:"mentioned_athletes"`
	PrimaryMediaType         string        `json:"primary_media_type"`
	ActivityType             string        `json:"activity_type"`
	IsFollowingActivityOwner bool          `json:"is_following_activity_owner"`
	StartDate                time.Time     `json:"start_date"`
	ElapsedTime              float32       `json:"elapsed_time"`
	Source                   string        `json:"source"`
	Rank                     int           `json:"rank"`
}

type Child struct {
	// Always appears to be empty?
}

type Submodule struct {
	Type                   string               `json:"type"`
	ModuleFields           []ModuleField        `json:"module_fields"`
	Destination            *Destination         `json:"destination,omitempty"`
	Element                string               `json:"element,omitempty"`
	ShouldTrackImpressions bool                 `json:"should_track_impressions,omitempty"`
	TrackableID            string               `json:"trackable_id,omitempty"`
	AnalyticsProperties    *AnalyticsProperties `json:"analytics_properties,omitempty"`
}
