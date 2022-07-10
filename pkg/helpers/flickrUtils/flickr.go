package flickrUtils

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/masci/flickr.v2"
	"gopkg.in/masci/flickr.v2/photosets"
)

func Init(key, secret string) (*flickr.FlickrClient, error) {
	client := flickr.NewFlickrClient(key, secret)
	requestTok, _ := flickr.GetRequestToken(client)
	url, _ := flickr.GetAuthorizeUrl(client, requestTok)

	var authCode string
	prompt := &survey.Input{
		Message: fmt.Sprintf("Please login to Flickr at this URL: %s\n Then enter the string here:", url),
		Default: "",
	}
	err := survey.AskOne(prompt, &authCode, survey.WithValidator(survey.Required))
	if err != nil {
		return nil, err
	}

	accessTok, err := flickr.GetAccessToken(client, requestTok, authCode)
	if err != nil {
		return nil, err
	}
	client.OAuthToken = accessTok.OAuthToken
	client.OAuthTokenSecret = accessTok.OAuthTokenSecret
	return client, nil
}

func UploadPhoto(client *flickr.FlickrClient, src string) (string, error) {
	resp, err := flickr.UploadFile(client, src, flickr.NewUploadParams())
	if err != nil {
		return "", err
	}
	if resp.Status != "ok" {
		return "", fmt.Errorf(resp.Error.Message)
	}

	return resp.ID, nil
}

func CreatePhotoSet(client *flickr.FlickrClient, title, description, primaryPhotoID string) (string, error) {
	resp, err := photosets.Create(client, title, description, primaryPhotoID)
	if err != nil {
		return "", err
	}
	if resp.Status != "ok" {
		return "", errors.New(resp.Error.Message)
	}
	return resp.Set.Id, nil
}
