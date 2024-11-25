# Frameo Control API

This Go application provides an HTTP API to control a Frameo digital picture frame using ADB commands. It allows you to toggle the screen, adjust brightness, and navigate between images.

## Prerequisites

- Ensure ADB is installed and accessible from your command line.
- Connect your Frameo device to the computer via USB and verify ADB connection.

## Usage

1. Run the server:

   ```bash
   go run frameo.go
   ```

2. The server will start on port 5000. The following endpoints are available:

   - **GET /brightness**: Retrieves the current brightness level (0-255).
   
   - **POST /brightness**: Sets screen brightness. Send a raw integer in the request body.

   - **GET /screen**: Retrieves current screen state (`true` for awake, `false` for asleep).

   - **POST /screen**: Toggles screen state. Send either a boolean value or leave the body empty to toggle without specifying state.

   - **POST /next**: Swipes right to show the next image.

   - **POST /previous**: Swipes left to show the previous image.