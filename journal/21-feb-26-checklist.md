API Endpoints

**User Management**
POST /signup: Register a new user.

POST /login: Authenticate a user.

GET /profile: Get user profile details.

PUT /profile: Update user profile.

**Blog Post Management**
POST /posts: Create a new blog post.

GET /posts: Retrieve a list of blog posts.

GET /posts/{id}: Retrieve a single blog post by ID.

PUT /posts/{id}: Update a blog post by ID.

DELETE /posts/{id}: Delete a blog post by ID.

**Comments**
POST /posts/{id}/comments: Add a comment to a blog post.

GET /posts/{id}/comments: Retrieve comments for a blog post.

DELETE /comments/{id}: Delete a comment by ID.

**Follow System**
POST /follow/{userId}: Follow a user.

DELETE /unfollow/{userId}: Unfollow a user.

GET /followers: Retrieve followers of the authenticated user.

GET /following: Retrieve users followed by the authenticated user.