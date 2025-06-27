package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.POST("/upload", rt.wrap(rt.uploadImage))

	rt.router.PUT("/me/username", rt.wrap(rt.idVerifierMiddleware(rt.setMyUsername)))
	rt.router.PUT("/me/photo", rt.wrap(rt.idVerifierMiddleware(rt.setMyPhoto)))

	rt.router.POST("/conversations", rt.wrap(rt.idVerifierMiddleware(rt.createConversation)))
	rt.router.GET("/conversations", rt.wrap(rt.idVerifierMiddleware(rt.getMyConversations)))
	rt.router.GET("/conversations/:conversationId", rt.wrap(rt.idVerifierMiddleware(rt.getConversation)))

	rt.router.PUT("/conversations/:conversationId/name", rt.wrap(rt.idVerifierMiddleware(rt.setGroupName)))
	rt.router.PUT("/conversations/:conversationId/photo", rt.wrap(rt.idVerifierMiddleware(rt.setGroupPhoto)))
	rt.router.POST("/conversations/:conversationId/participants", rt.wrap(rt.idVerifierMiddleware(rt.addToGroup)))
	rt.router.DELETE("/conversations/:conversationId/participants", rt.wrap(rt.idVerifierMiddleware(rt.leaveGroup)))

	rt.router.POST("/conversations/:conversationId/messages", rt.wrap(rt.idVerifierMiddleware(rt.sendMessage)))
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId", rt.wrap(rt.idVerifierMiddleware(rt.deleteMessage)))
	rt.router.POST("/conversations/:conversationId/forwarded_messages", rt.wrap(rt.idVerifierMiddleware(rt.forwardMessage)))

	rt.router.POST("/conversations/:conversationId/messages/:messageId/reactions", rt.wrap(rt.idVerifierMiddleware(rt.commentMessage)))
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId/reactions", rt.wrap(rt.idVerifierMiddleware(rt.uncommentMessage)))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
