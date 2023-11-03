# RPC Transport for refactored-robot

This package defines a gRPC service for user management.

## Protocol Buffers (protobuf)

The service uses Protocol Buffers (protobuf) with syntax version "proto3".

## Go Package

The Go package for this service is located in the `./pb` directory.

## Service

### UserController

The `UserController` service provides the following RPC methods:

- `AddUser(RegisterUserRequest) returns (google.protobuf.Empty)`: Adds a new user.
- `GetUser(GetUserRequest) returns (UserResponse)`: Retrieves user information.
- `DeleteUser(DeleteUserRequest) returns (google.protobuf.Empty)`: Deletes a user.
- `Login(LoginRequest) returns (LoginResponse)`: Performs user login and returns authentication tokens.
- `UploadImage(UploadImageRequest) returns (google.protobuf.Empty)`: Uploads an image associated with a user.
- `GetImage(GetImageRequest) returns (GetImageResponse)`: Retrieves user's image data.
- `RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse)`: Refreshes authentication tokens.

## Message Types

### User

- `id` (int32): User ID.
- `name` (string): User name.
- `password` (string): User password.
- `image` (string): URL or path to user's image.

### RegisterUserRequest

- `name` (string): User name for registration.
- `password` (string): User password for registration.
- `image` (string): URL or path to user's image for registration.

### UserResponse

- `id` (int32): User ID.
- `name` (string): User name.
- `password` (string): User password.
- `image` (string): URL or path to user's image.

### GetUserRequest

- `id` (int32): User ID for retrieval.

### DeleteUserRequest

- `id` (int32): User ID for deletion.

### LoginRequest

- `name` (string): User name for login.
- `password` (string): User password for login.

### LoginResponse

- `token` (string): Authentication token.
- `refreshToken` (string): Refresh token for authentication.

### UploadImageRequest

- `id` (int32): User ID associated with the image.
- `imageBytes` (bytes): Image data in bytes.

### GetImageRequest

- `id` (int32): User ID for image retrieval.

### GetImageResponse

- `imageData` (bytes): Image data in bytes.

### RefreshTokenRequest

- `refreshToken` (string): Refresh token for token refresh.

### RefreshTokenResponse

- `token` (string): New authentication token.

### Success

- `Success` (string): A success message.

## Examples

### AddUser

```protobuf
rpc_transport.UserController.AddUser({
  name: "John Doe",
  password: "password123",
  image: "/images/john_doe.jpg"
});
```
### GetUser
```protobuf
rpc_transport.UserController.GetUser({
  id: 123
});
```
### DeleteUser
```protobuf
rpc_transport.UserController.DeleteUser({
  id: 123
});
```
### Login
```protobuf
rpc_transport.UserController.Login({
  name: "John Doe",
  password: "password123"
});
```
### UploadImage
```protobuf
rpc_transport.UserController.UploadImage({
  id: 123,
  imageBytes: <image_data_bytes>
});
```
### GetImage
```protobuf
rpc_transport.UserController.GetImage({
  id: 123
});
```
### RefreshToken
```protobuf
rpc_transport.UserController.RefreshToken({
  refreshToken: "refresh_token_string"
});
```