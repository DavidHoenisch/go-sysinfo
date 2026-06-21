package gosysinfo

import (
	"reflect"
	"testing"
)

func TestGetScreenVirtualSize(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    ScreenSize
	}{
		{
			name:    "valid input",
			content: "1920,1080",
			want:    ScreenSize{X: "1920", Y: "1080"},
		},
		{
			name:    "zero width",
			content: "0,0",
			want:    ScreenSize{X: "0", Y: "0"},
		},
		{
			name:    "empty content",
			content: "",
			want:    ScreenSize{X: "0", Y: "0"},
		},
		{
			name: "nil content",
			want: ScreenSize{X: "0", Y: "0"},
		},
		{
			name:    "invalid input",
			content: "invalid",
			want:    ScreenSize{X: "0", Y: "0"},
		},
		{
			name:    "large width",
			content: "3840,2160",
			want:    ScreenSize{X: "3840", Y: "2160"},
		},
		{
			name:    "small height",
			content: "1920,720",
			want:    ScreenSize{X: "1920", Y: "720"},
		},
		{
			name:    "zero height",
			content: "1920,0",
			want:    ScreenSize{X: "1920", Y: "0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := mockReader{}
			if tt.name != "nil content" {
				reader[fb0_virtual_size] = tt.content
			}

			if got := GetScreenVirtualSize(reader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetScreenVirtualSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetScreenMode(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    ScreenMode
	}{
		{
			name:    "valid input",
			content: "1920x1080",
			want:    ScreenMode{Value: "1920x1080"},
		},
		{
			name:    "multiline modes",
			content: "1920x1080\n1280x720",
			want:    ScreenMode{Value: "1920x1080\n1280x720"},
		},
		{
			name:    "empty content",
			content: "",
			want:    ScreenMode{Value: ""},
		},
		{
			name: "nil content",
			want: ScreenMode{Value: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := mockReader{}
			if tt.name != "nil content" {
				reader[fb0_screen_modes] = tt.content
			}

			if got := GetScreenMode(reader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetScreenMode() = %v, want %v", got, tt.want)
			}
		})
	}
}
