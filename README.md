# Forum

## Created by:
1. nsergaziy (Nurkhat)
2. drakhmat (Dayana)

## Overview
The Forum project is an interactive, user-friendly platform designed for discussing and sharing ideas. It allows users to post content, comment, like posts, and manage profiles while providing a seamless experience with detailed filters and robust error handling. Additionally, it incorporates an admin panel for user management and monitoring.

The project follows best practices in interface design, adhering to Schneiderman's 8 Golden Rules, ensuring consistency, responsiveness, and accessibility. The backend is powered by Go, with Docker support for easy deployment.


## Project Structure

The Groupie Tracker Filters project is structured as follows:

```
/forum
├── /cmd
│   └── main.go
├── /internal
│   ├── /database
│   │   ├── dummy.db
│   │   └── init.sql
│   ├── /handlers 
│   │   ├── comment.go
│   │   ├── errors.go
│   │   ├── home.go
│   │   ├── main_test.go
│   │   ├── post.go
│   │   ├── user.go
│   │   └── vote.go
│   ├── /models
│   │   ├── comment.go
│   │   ├── post.go
│   │   └── user.go
│   └── routes.go
├── /ui
│   ├── /static
│   │   ├── /css
│   │   │   └── styles.css
│   │   ├── /img
│   │   │   ├── bow.png
│   │   │   ├── comment.png
│   │   │   ├── disike.png
│   │   │   ├── error.png
│   │   │   ├── bow.png
│   │   │   ├── like.png
│   │   │   └── otter.png
│   │   └── /js
│   │       └── main.js
│   └── /templates
│       ├── create.html
│       ├── error.html
│       ├── footer.html
│       ├── header.html
│       ├── home.html
│       ├── left_sidebar.html
│       ├── login.html
│       ├── profile.html
│       ├── right_sidebar.html
│       ├── signup.html
│       └── view.html
├──  .dockerignore
├──  docker-compose.yml
├──  DockerFile
├──  go.mod
└──  README.md
```

## Objective

The Forum project aims to:
1. Provide a robust platform for users to interact via posts, comments, and likes. 
2. Ensure scalability and ease of deployment using Docker. 
3. Incorporate advanced filters for personalized user experiences. 
4. Offer an admin panel for managing users and monitoring activities. 
5. Follow best practices in UI/UX design for a consistent and responsive layout.

## Docker and Deployment
The project includes Docker support for easy setup and deployment:
1. Dockerfile
   - Builds a containerized Go application.
2. docker-compose.yml:
   - Configures a multi-service setup, including the application server and database.

## Usage
1. Clone the repository: `git@git.01.alem.school:nsergaziy/forum.git`.
2. Navigate to the project directory: `cd forum`.
3. Build and run using Docker:
    - To start the server: `docker-compose up --build`.
4. Open your web browser and navigate to `http://localhost:8080`.

## Features

### User Management
1. Registration & Authentication:
   - Secure user registration and login using hashed passwords.
   - Passwords are stored securely with bcrypt.
   - Persistent login using session cookies.
2. Profile Management:
   - Users can update their profile information, including username and password.
   - A detailed user dashboard showcasing personal posts, liked posts, and comments.
### Post Interactions
1. Posting:
   - Users can create posts and choose Category.
   - View a list of posts created by user (My Posts).
2. Likes:
   - Like and dislike posts.
   - View a list of liked posts (Liked Posts)
3. Comments:
- Add comments to posts.
- View all personal comments in the (Commented Posts).
### Category Filters
The application includes powerful category filters for posts:
- Technology
- Entertainment
- Sports
- Education
- Health
### Admin Panel
1. Default admin credentials:
   - Email: admin@gmail.com
   - Password: 12345678
2. Admin Controls:
   - View all registered users.
   - Ban or unban users.
   - Monitor posts and comments for inappropriate content.


## Testing

### Key Testing Areas
The handlers package in the Forum application includes a series of tests to ensure that API endpoints and database interactions function correctly. This testing suite uses httptest to simulate HTTP requests and assert from the testify package for validation. Mock objects are employed to isolate the tests and mimic database behavior.
1. GetHandler Tests:
   - Verifies the GetHandler correctly fetches data and handles errors.
   - Ensures proper HTTP status codes and headers are set.
2. PostHandler Tests:
   - Validates the insertion of new data through the PostHandler.
   - Tests edge cases, such as empty body submissions.
3. AddPostToCategoryHandler Tests:
   - Confirms correct linking of posts to categories, ensuring post IDs and category IDs are passed correctly.
4. GetCategory Tests:
   - Validates fetching of category data by ID, including handling errors like "not found."
5. Stress Testing:
   - Simulates multiple concurrent requests to ensure the handlers can handle high traffic without issues.

### Running Tests
To run the tests, use the following command:
```
go test ./internal/handlers
```

## Conclusion
The Forum project is a comprehensive and interactive platform for engaging discussions. With powerful features, intuitive design, and scalable architecture, it serves as a good project for interaction. Hope to your audit.