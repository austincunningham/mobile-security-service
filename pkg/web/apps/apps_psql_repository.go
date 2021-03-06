package apps

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	"github.com/aerogear/mobile-security-service/pkg/models"
)

type (
	// PostgreSQLRepository interface defines the methods to be implemented
	appsPostgreSQLRepository struct {
		db *sql.DB
	}
)

// NewPostgreSQLRepository creates a new instance of appsPostgreSQLRepository
func NewPostgreSQLRepository(db *sql.DB) Repository {
	return &appsPostgreSQLRepository{db}
}

// GetApps retrieves all apps from the database
func (a *appsPostgreSQLRepository) GetApps() (*[]models.App, error) {
	rows, err := a.db.Query(`
	SELECT a.id,a.app_id,a.app_name,
	COALESCE(COUNT(DISTINCT v.id),0) as num_of_deployed_versions,
	COALESCE(SUM(DISTINCT v.num_of_app_launches),0) as num_of_app_launches,
	COALESCE(COUNT(DISTINCT d.id),0) as num_of_current_installs
	FROM app as a LEFT JOIN version as v on a.app_id = v.app_id 
	LEFT JOIN device as d on v.id = d.version_id 
	WHERE a.deleted_at IS NULL 
	GROUP BY a.id;`)

	if err != nil {
		log.Error(err)
		return nil, models.ErrInternalServerError
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	apps := []models.App{}
	for rows.Next() {
		var a models.App
		if err = rows.Scan(&a.ID, &a.AppID, &a.AppName, &a.NumOfDeployedVersions, &a.NumOfAppLaunches, &a.NumOfCurrentInstalls); err != nil {
			log.Error(err)
		}

		apps = append(apps, a)
	}

	if len(apps) == 0 {
		return nil, models.ErrNotFound
	}

	return &apps, nil
}

func (a *appsPostgreSQLRepository) GetAppVersionsByAppID(id string) (*[]models.Version, error) {
	rows, err := a.db.Query(`
	SELECT v.id,v.version,v.id, v.disabled, v.disabled_message, v.num_of_app_launches,
	COALESCE(COUNT(DISTINCT d.id),0) as num_of_current_installs
	FROM version as v LEFT JOIN device as d on v.id = d.version_id
	WHERE v.app_id = $1 
	GROUP BY v.id;`, id)

	if err != nil {
		log.Error(err)
		return nil, models.ErrInternalServerError
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	versions := []models.Version{}

	// iterate over the rows and add the data to the array of versions
	for rows.Next() {
		var v models.Version
		var disabledMessage sql.NullString
		if err = rows.Scan(&v.ID, &v.Version, &v.AppID, &v.Disabled, &disabledMessage, &v.NumOfCurrentInstalls, &v.NumOfAppLaunches); err != nil {
			log.Error(err)
		}

		v.DisabledMessage = disabledMessage.String
		versions = append(versions, v)
	}

	if len(versions) == 0 {
		return nil, models.ErrNotFound
	}

	return &versions, nil
}
