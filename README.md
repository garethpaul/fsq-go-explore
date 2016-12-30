Foursquare-Go-Explore
===================
An example of how to access Foursquare Data using GoLang and AppEngine.

Getting Started
----------
Here are some steps to get you started.

#### 1. Get your developer keys

Sign up for a <a href="https://developer.foursquare.com">Foursquare Developer</a> or <a href="https://enterprise.foursquare.com/contact-us">Enterprise Account</a>.

Once you get setup you should  have two important pieces of information:

1. CLIENT_ID - Unique to your registered application.
2. CLIENT_SECRET - Unique and private to your application.

Please remember to keep these keys in a secure location.


#### 2. Clone this repo

``` git clone https://github.com/garethpaul/fsq-go-explore.git```

#### 3. Amend App.yaml with your variables

These are directly environment variables.

```
env_variables:
  FSQ_CLIENT_ID: 'YOUR_FOURSQUARE_KEY' // found in step 1
  FSQ_CLIENT_SECRET: 'YOUR_CLIENT_SECRET' // found in step 1
  FSQ_VERSION: 'YYYYMMDD' // e.g. 20170101
```

##### 4. Run your application

```
goapp serve
```

> **Note:**

> - You'll need to install [GoLang](https://golang.org/doc/install)
> - You will need to download the [Google AppEngine SDK for GoLang](https://cloud.google.com/appengine/docs/go/download)
