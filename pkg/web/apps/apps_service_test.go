package apps

import (
	"reflect"
	"testing"

	"github.com/aerogear/mobile-security-service/pkg/models"
)

func Test_appsService_GetApps(t *testing.T) {
	tests := []struct {
		name    string
		want    *[]models.App
		wantErr bool
	}{
		{
			name: "appsService.GetApps() should return a list of apps",
			want: &[]models.App{
				models.App{
					ID:                    1,
					AppID:                 "com.aerogear.app1",
					AppName:               "app1",
					NumOfDeployedVersions: 1,
					NumOfAppLaunches:      1,
					NumOfClients:          1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		appsPostgreSQLRepository := NewPostgreSQLRepository()

		t.Run(tt.name, func(t *testing.T) {
			a := &appsService{
				repository: appsPostgreSQLRepository,
			}

			got, err := a.GetApps()
			if (err != nil) != tt.wantErr {
				t.Errorf("appsService.GetApps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appsService.GetApps() = %v, want %v", got, tt.want)
			}
		})
	}
}
