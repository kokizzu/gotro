
# Tutorial

## Google OAuth

1. open https://console.cloud.google.com/apis/credentials?project=
2. create new project, fill OAuth Consent Screen requirements, tick (openid,email,profile), create credentials "OAuth 2.0 CLient IDs" > Web Application
3. fill Authorized Javascript origins with `http://localhost:9090` and `http://127.0.0.1:9090`
4. fill Authorized Redirect URIs with `http://localhost:9090/api/UserOauth` and `http://127.0.0.1:9090/api/UserOauth`
5. put the ClientID and ClientSecret to `.env`

## Yahoo

doc: https://developer.yahoo.com/oauth2/guide/flows_authcode/

1. open https://developer.yahoo.com/apps/
2. create an app, tick the OpenID integration, fill your appname, description, homepage url and redirect url
3. fill homepage url with `http://127.0.0.1`
4. fill redirect uri with `https://127.0.0.1/api/UserOauth` and `https://localhost/api/UserOauth`
5. put the ClientID and ClientSecret to `.env`

note:
- Yahoo doesn't support http Redirect URIs, so use Caddy to run https locally
- yahoo does not support multiple origin domain nor IP, so we have to create one by one

## Facebook

TODO

## Github

1. Login to your github account
2. Go to settings > developer setting > oauth apps > new apps
3. Fill your app name, homepage url, description, auth callback url.
4. fill homepage url with `http://127.0.0.1`
5. fill callback url with `http://localhost/api/UserOauth`
6. Click generate `Client Secret` and don't forget to copy
7. Submit `Update Application`
8. put the ClientID and ClientSecret to `.env`

note:
- If you forgot copy client secret, you can delete first older cilent secret and generate again, then copy, submit app
- Github only allow add 1 callback url per app

## Twitter

TODO

## Steam

TODO

# How to test

1. start the apiserver (`make apiserver`), then hit one of these url to retrieve the login link
[http://localhost:9090/api/UserExternalLogin?provider=google](http://localhost:9090/api/UserExternalLogin?provider=google)
[http://localhost:9090/api/UserExternalLogin?provider=yahoo](http://localhost:9090/api/UserExternalLogin?provider=yahoo)

2. login with your gmail/yahoo, then it would give a response like this:

```json
{"sessionToken":"0QCPcCddiD-~-----2~0|1d372ee5601c427ec97df10fa2bee413660d952e7a0c57d24bba0bf31b50644ca2a2dcb216758c28b488dd4c86a47a41e988c6efbe2a43bcd0794926|0T0mT4j4Koq","error":"","status":0,"OauthUser":{"email":"xxx@gmail.com","email_verified":true,"family_name":"xxx","given_name":"xxx","locale":"en-US","name":"xxx xxx","picture":"https://lh3.googleusercontent.com/a-/AOh14GjfGEGb0xaIbnjdRnZWas3NBhoYdQaCcFb66Pbcag=s96-c","sub":"101959137763910089936"},"Email":"xxx@gmail.com","currentUser":{"id":"144428372796112898","email":"xxx@gmail.com","password":"$2a$10$onkM0VBO90l3DBiRh2sqTeegZKE3JcKIWxzS3clb5rKDI.kjQAIqC","createdAt":0,"createdBy":"0","updatedAt":0,"updatedBy":"0","deletedAt":0,"deletedBy":"0","isDeleted":false,"restoredAt":0,"restoredBy":"0","passwordSetAt":0,"secretCode":"","secretCodeAt":0,"verificationSentAt":0,"verifiedAt":0,"lastLoginAt":0}}
```

```json
{"sessionToken":"0QCUkl01IP7~-----1~0|8c72df957ed26b3aefe1c9b46753769587247af21583bccfe723facde8c7cef88311f43ad3e7c711ecbf3daa797313d24474afe41a8545961ffbb54a|0T0mT4j4Koq","error":"","status":0,"oauthUser":{"email":"xxx@gmail.com","email_verified":false,"family_name":"xxx","gender":"notDisclosed","given_name":"xxx","locale":"id-ID","name":"xxx xxx","nickname":"xxx","picture":"https://s.yimg.com/ag/images/default_user_profile_pic_192sq.jpg","profile_images":{"image128":"https://s.yimg.com/ag/images/default_user_profile_pic_128sq.jpg","image192":"https://s.yimg.com/ag/images/default_user_profile_pic_192sq.jpg","image32":"https://s.yimg.com/ag/images/default_user_profile_pic_32sq.jpg","image64":"https://s.yimg.com/ag/images/default_user_profile_pic_64sq.jpg"},"sub":"3ZQ6ACERLUCI5WMDFQYNHYIRKU"},"email":"xxx@gmail.com","currentUser":{"id":"144428372796112898","email":"xxx@gmail.com","password":"$2a$10$onkM0VBO90l3DBiRh2sqTeegZKE3JcKIWxzS3clb5rKDI.kjQAIqC","createdAt":0,"createdBy":"0","updatedAt":0,"updatedBy":"0","deletedAt":0,"deletedBy":"0","isDeleted":false,"restoredAt":0,"restoredBy":"0","passwordSetAt":0,"secretCode":"","secretCodeAt":0,"verificationSentAt":0,"verifiedAt":0,"lastLoginAt":0}}
```
