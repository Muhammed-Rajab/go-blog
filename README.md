# go-blog

_✨A minimalistic blog made with zero distractions✨_

## Run Locally

Requirements

```
go >= 1.19
mongodb running locally
```

Clone the project

```bash
  git clone https://github.com/Muhammed-Rajab/go-blog
```

Go to the project directory

```bash
  cd go-blog/
```

Install dependencies

```bash
  go get
```

Create a .env file

```bash
  touch .env
```

Copy the following to the .env file

```bash
  BLOG_DASHBOARD_KEY=<YOUR_DASHBOARD_SECRET>
  MONGODB_URI=<YOUR_MONGODB_URI>
```

Start the server

```bash
  go run .
```

The blog will start serving at http://localhost:8000/blog

## Documentation

### Dashboard

#### Authentication

**go-blog** employs a straightforward authentication mechanism. Set the `BLOG_DASHBOARD_KEY` as an environment variable, serving as your access token to the **dashboard**.

To access the dashboard, visit [http://localhost:8000/blog/dashboard/auth](http://localhost:8000/blog/dashboard/auth) and enter your `BLOG_DASHBOARD_KEY`. This action will establish a cookie in your browser, allowing the server to identify you. Feel free to change the `BLOG_DASHBOARD_KEY` whenever needed, and remember to restart the server afterward.

#### Home

The dashboard home is designed to be intuitive and user-friendly. Take a moment to explore its features and discover everything it has to offer on your own. It's designed to be self-explanatory, making your navigation seamless and enjoyable. Happy exploring!

#### Uploading Images

At present, **go-blog** doesn't include a built-in markdown editor. You can upload images to the server, provided they are under the size limit of `10MB`. The uploaded images can be viewed at [http://localhost:8000/blog/dashboard/images](http://localhost:8000/blog/dashboard/images). Simply click on an uploaded image to obtain its link, which can then be added to your article.

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`AUTHOR`: blog author's name

`BLOG_DASHBOARD_KEY`: your secret to access the dashboard

`UPLOADS_DIR`: path to store the uploaded images

`MONGODB_URI`: your mongodb uri

## Screenshots

![Home 1](/screenshots/home-1.png)
![Home 2](/screenshots/home-2.png)
![Blog 1](/screenshots/blog-1.png)
![Blog 2](/screenshots/blog-2.png)
![Blog 2](/screenshots/blog-2.png)
![Dashboard 1](/screenshots/dashboard-1.png)
![Dashboard 2](/screenshots/dashboard-2.png)
![Add blog](/screenshots/add-blog.png)
![Images 1](/screenshots/images-1.png)
![Images 2](/screenshots/images-2.png)
![Images 3](/screenshots/images-3.png)

## License

[MIT](https://choosealicense.com/licenses/mit/)
