
# Tutorial

## Google OAuth

1. open https://console.cloud.google.com/apis/credentials?project=
2. create new project, fill OAuth Consent Screen requirements, tick (openid,email,profile), create credentials "OAuth 2.0 CLient IDs" > Web Application
3. fill Authorized Javascript origins with `http://localhost:9090` and `http://127.0.0.1:9090`
4. fill Authorized Redirect URIs with `http://localhost:9090/api/UserOauth` and `http://127.0.0.1:9090/api/UserOauth`

## Facebook

TODO

## Github

TODO

## Twitter

TODO

## Yahoo

NOTE: yahoo does not support multiple domain, does not support IP, so we have to create one by one.

1. 

## Steam

TODO

# How to test

1. start the apiserver (`make apiserver`), then hit the url to retrieve the login link
[http://localhost:9090/api/UserExternalLogin?provider=google](http://localhost:9090/api/UserExternalLogin?provider=google)


2. do login with your gmail, then it would give a response like this:

```json
{"sessionToken":"0QCPcCddiD-~-----2~0|1d372ee5601c427ec97df10fa2bee413660d952e7a0c57d24bba0bf31b50644ca2a2dcb216758c28b488dd4c86a47a41e988c6efbe2a43bcd0794926|0T0mT4j4Koq","error":"","status":0,"OauthUser":{"email":"xxx@gmail.com","email_verified":true,"family_name":"xxx","given_name":"xxx","locale":"en-US","name":"xxx xxx","picture":"https://lh3.googleusercontent.com/a-/AOh14GjfGEGb0xaIbnjdRnZWas3NBhoYdQaCcFb66Pbcag=s96-c","sub":"101959137763910089936"},"Email":"xxx@gmail.com","CurrentUser":{"id":"144428372796112898","email":"xxx@gmail.com","password":"$2a$10$onkM0VBO90l3DBiRh2sqTeegZKE3JcKIWxzS3clb5rKDI.kjQAIqC","createdAt":0,"createdBy":"0","updatedAt":0,"updatedBy":"0","deletedAt":0,"deletedBy":"0","isDeleted":false,"restoredAt":0,"restoredBy":"0","passwordSetAt":0,"secretCode":"","secretCodeAt":0,"verificationSentAt":0,"verifiedAt":0,"lastLoginAt":0}}
```
