/*
 * Script Type: OAuth2 Access Token Modification
 * Purpose: Inject jwks_uri from client properties into the Access Token
 */
(function () {
    // Extract the custom jwks_uri property
    // (Assumes it was saved as 'jwks_uri' during Dynamic Registration)
    var jwksUri = clientProperties.get('customProperties').get("jwks_uri")
    if (jwksUri != null && !jwksUri.isEmpty()) {
        // Add the claim to the Access Token
        accessToken.setField("jwks_uri", jwksUri)
    } else {
        // Optional: throw an exception to prevent invalid token generation
        logger.warning("No jwks_uri found in configuration for client: " + clientProperties.get('clientId'))
    }
    // No return value is expected. Let it be undefined.
}());