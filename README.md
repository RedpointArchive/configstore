# configstore

*configstore is a Go server that lets you define your schema in JSON, and automatically generates gRPC calls to List/Get/Create/Update/Delete those entities in Google Cloud Firestore.*

Basically, configstore is a server that provides a gRPC CRUD interface for the kinds of entities that you describe. You just specify the schema for kinds as JSON, and configstore produces a `.proto` file you can use to generate clients in various languages.

configstore is entirely stateless; all data is stored on Google Cloud Firestore. This makes it easier to run in highly available configurations, without managing the state yourself.

## Usage

configstore is a very early prototype, and isn't ready for production usage.

You can run the test server on Windows by running `.\server.ps1`. You'll most likely need to authenticate your local computer with Application Default Credentials, and set `$env:GOOGLE_CLOUD_PROJECT_ID` to the Google Cloud project that is running your Firestore instance.

You can test that it works for the example schema by running `.\client.ps1`. You'll need to manually create a `User` entity in Firestore with an `emailAddress` string property, since the only implemented gRPC method for kinds in configstore is `Get`.

## Screenshots

![server](./screenshots/server.png)

![client](./screenshots/client.png)
