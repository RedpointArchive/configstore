# configstore

*configstore is a Go server that lets you define your schema in JSON, and automatically generates gRPC calls to List/Get/Create/Update/Delete those entities in Google Cloud Firestore.*

Basically, configstore is a server that provides a gRPC CRUD interface for the kinds of entities that you describe. You just specify the schema for kinds as JSON, and configstore produces a `.proto` file you can use to generate clients in various languages.

configstore is entirely stateless; all data is stored on Google Cloud Firestore. This makes it easier to run in highly available configurations, without managing the state yourself.

## Usage

configstore is a very early prototype, and isn't ready for production usage.

You can run the test server on Windows by running `.\server.ps1`. You'll most likely need to authenticate your local computer with Application Default Credentials, and set `$env:GOOGLE_CLOUD_PROJECT_ID` to the Google Cloud project that is running your Firestore instance.

You can test that it works for the example schema by running `.\client.ps1`. You'll need to manually create a `User` entity in Firestore with an `emailAddress` string property, since the only implemented gRPC method for kinds in configstore is `Get`.

## Screenshots

![server](https://github.com/hach-que/configstore/raw/master/screenshots/server.PNG)

![client](https://github.com/hach-que/configstore/raw/master/screenshots/client.PNG)

## License

```
MIT License

Copyright (c) 2019 June Rhodes

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```