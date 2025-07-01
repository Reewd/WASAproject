<template>
	<div
		:class="['message-container', { 'own-message': isOwnMessage }]"
		@contextmenu.prevent="showContextMenu"
		@click="hideContextMenu"
	>
		<!-- Existing message content -->
		<div v-if="isGroupConversation && !isOwnMessage" class="message-header">
			<img
				:src="
					getImageUrl(message.sentBy.photo?.path) || userDefaultIcon
				"
				alt="User Photo"
				class="profile-picture"
			/>
			<span class="username">{{ message.sentBy.username }}</span>
		</div>

		<div class="message-wrapper">
			<div :class="['message-bubble', { 'own-bubble': isOwnMessage }]">
				<!-- Reaction Button (shows on hover) -->
				<ReactionButton
					:message="message"
					:isOwnMessage="isOwnMessage"
					:conversationId="conversationId"
					@openEmojiPicker="handleOpenEmojiPicker"
				/>

				<!-- Forwarded Message Indicator -->
				<div v-if="message.isForwarded" class="forwarded-indicator">
					<span class="forwarded-icon">↪️</span>
					<span class="forwarded-text">Forwarded</span>
				</div>

				<!-- Reply to message -->
				<div v-if="message.replyTo" class="reply-container">
					<div class="reply-header">
						<strong>{{
							replyToMessage?.sentBy.username || "Unknown User"
						}}</strong>
					</div>
					<div class="reply-content">
						<div v-if="replyToMessage?.photo" class="reply-with-photo">
							<img
								:src="getImageUrl(replyToMessage.photo.path)"
								alt="Reply Photo"
								class="reply-photo"
							/>
							<div class="reply-text-content">
								<span v-if="replyToMessage?.text">{{
									replyToMessage.text
								}}</span>
								<span v-else class="attachment-indicator">📷 Photo</span>
							</div>
						</div>
						<span v-else-if="replyToMessage?.text">{{
							replyToMessage.text
						}}</span>
						<span v-else>Message not available</span>
					</div>
				</div>

				<!-- Message content -->
				<div class="message-content">
					<!-- Photo -->
					<img
						v-if="message.photo"
						:src="getImageUrl(message.photo.path)"
						alt="Message Photo"
						class="message-photo"
					/>

					<!-- Text -->
					<p v-if="message.text" class="message-text">
						{{ message.text }}
					</p>
				</div>

				<!-- Message metadata -->
				<div class="message-metadata">
					<span class="timestamp">{{
						formatTimestamp(message.timestamp)
					}}</span>
					<span
						v-if="isOwnMessage"
						:class="['status', message.status]"
					>
						{{ getStatusIcon(message.status) }}
					</span>
				</div>
			</div>

			<!-- Reaction Bubble -->
			<ReactionBubble
				v-if="message.reactions && message.reactions.length > 0"
				:reactions="message.reactions"
				:class="[
					'reaction-bubble-container',
					{ 'own-reactions': isOwnMessage },
				]"
				@removeReaction="handleRemoveReaction"
			/>
		</div>

		<!-- Context Menu -->
		<MessageContextMenu
			:isVisible="showMenu"
			:position="menuPosition"
			:message="message"
			@close="hideContextMenu"
			@reply="handleReply"
			@copy="handleCopy"
			@forward="handleForward"
			@delete="handleDelete"
		/>

		<!-- Forward Message Modal -->
		<ForwardMessage
			v-if="showForwardModal"
			:message="message"
			:currentConversationId="conversationId"
			@close="closeForwardModal"
			@forwarded="handleMessageForwarded"
		/>
	</div>
</template>

<script setup>
import { ref, computed } from "vue";
import axios from "../services/axios.js";
import { useUser } from "../composables/useUser.js";
import { useImageUrl } from "../composables/useImageUrl.js";
import ReactionButton from "./ReactionButton.vue";
import MessageContextMenu from "./MessageContextMenu.vue";
import ReactionBubble from "./ReactionBubble.vue";
import ForwardMessage from "../modals/ForwardMessage.vue"; // Add this import
import userDefaultIcon from "/assets/icons/user-default.png";

const { getImageUrl } = useImageUrl();
const { getUsername } = useUser();
const { getUserId } = useUser();
const currentUsername = getUsername();
const showForwardModal = ref(false);
const handleForward = (message) => {
	console.log("Forward message:", message);
	showForwardModal.value = true;
};

const closeForwardModal = () => {
	showForwardModal.value = false;
};

const handleMessageForwarded = (data) => {
	console.log("Message forwarded:", data);
	showForwardModal.value = false;
	emits("messageForwarded");
};
const handleRemoveReaction = async () => {
	try {
		// Since the backend removes reactions based on messageId + userId,
		// we don't need to send reactionId or emoji in the request
		await axios.delete(
			`/conversations/${props.conversationId}/messages/${props.message.messageId}/reactions`,
			{
				headers: {
					Authorization: getUserId(),
				},
			}
		);

		emits("reactionRemoved");
	} catch (error) {
		console.error("Error removing reaction:", error);
		// Log the full error response to see what's causing the bad request
		if (error.response) {
			console.error("Error response:", error.response.data);
			console.error("Error status:", error.response.status);
		}
	}
};
const props = defineProps({
	message: {
		type: Object,
		required: true,
	},
	isGroupConversation: {
		type: Boolean,
		default: false,
	},
	replyToMessage: {
		type: Object,
		default: null,
	},
	conversationId: {
		type: Number,
		required: true,
	},
});

const emits = defineEmits([
	"reply",
	"reactionRemoved",
	"openEmojiPicker",
	"messageDeleted",
	"messageForwarded",
]);

// Context menu state
const showMenu = ref(false);
const menuPosition = ref({ x: 0, y: 0 });

// Check if this message is sent by the current user
const isOwnMessage = computed(() => {
	return props.message.sentBy.username === currentUsername;
});

// Context menu handlers
const showContextMenu = (event) => {
	event.preventDefault();
	menuPosition.value = {
		x: event.clientX,
		y: event.clientY,
	};
	showMenu.value = true;
};

const hideContextMenu = () => {
	showMenu.value = false;
};

const handleReply = (message) => {
	console.log("Reply to message:", message);
	emits("reply", message);
};

const handleCopy = (message) => {
	console.log("Copy message:", message);
	// Copy logic would go here
};

const handleDelete = async (message) => {
	try {
		await axios.delete(
			`/conversations/${props.conversationId}/messages/${props.message.messageId}`,
			{
				headers: {
					Authorization: getUserId(),
				},
			}
		);

		console.log("Message deleted successfully");
		// Emit event to parent component to handle UI updates
		emits("messageDeleted", message.messageId);
	} catch (error) {
		console.error("Error deleting message:", error);
		if (error.response) {
			console.error("Error response:", error.response.data);
			console.error("Error status:", error.response.status);

			if (error.response.status === 403) {
				alert("You can only delete your own messages.");
			} else if (error.response.status === 404) {
				alert("Message not found.");
			} else {
				alert("Failed to delete message. Please try again.");
			}
		} else {
			alert("Failed to delete message. Please check your connection.");
		}
	}
};

// Handle emoji picker open request
const handleOpenEmojiPicker = (data) => {
	emits("openEmojiPicker", {
		messageId: data.messageId,
		position: data.position,
		conversationId: props.conversationId,
	});
};

// Format timestamp
const formatTimestamp = (timestamp) => {
	const date = new Date(timestamp);
	return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
};

// Get status icon
const getStatusIcon = (status) => {
	switch (status) {
		case "sent":
			return "☑️";
		case "delivered":
			return "☑️";
		case "read":
			return "✅✅";
		default:
			return "";
	}
};
</script>

<style scoped>
.message-container {
	display: flex;
	flex-direction: column;
	margin: 8px 0;
	align-items: flex-start;
	position: relative;
}

.message-container.own-message {
	align-items: flex-end;
}

.message-header {
	display: flex;
	align-items: center;
	margin-bottom: 4px;
	margin-left: 8px;
}

.profile-picture {
	width: 24px;
	height: 24px;
	border-radius: 50%;
	margin-right: 8px;
}

.username {
	font-size: 12px;
	font-weight: bold;
	color: #666;
}

.message-wrapper {
	display: flex;
	flex-direction: column;
	align-items: flex-start;
	position: relative;
	width: 100%; /* Add this */
}

.message-container.own-message .message-wrapper {
	align-items: flex-end;
}

.message-bubble {
	min-width: 100px; /* Add minimum width */
	max-width: min(400px, 75%); /* Better max-width calculation */
	width: fit-content; /* Add this */
	background-color: #e5e5ea;
	border-radius: 18px;
	padding: 12px 16px;
	position: relative;
	transition: transform 0.1s ease;
	word-wrap: break-word; /* Ensure text wraps properly */
}

.message-bubble.own-bubble {
	background-color: #007aff;
	color: white;
}

.message-container:hover .message-bubble {
	transform: scale(1.02);
}

.reaction-bubble-container {
	margin-top: -8px;
	margin-left: 16px;
	z-index: 1;
}

.reaction-bubble-container.own-reactions {
	margin-left: 0;
	margin-right: 16px;
	align-self: flex-end;
}

.reply-container {
	background-color: rgba(0, 0, 0, 0.1);
	border-radius: 8px;
	padding: 8px;
	margin-bottom: 8px;
	border-left: 3px solid #007aff;
	width: 100%; /* Ensure reply container takes full bubble width */
	box-sizing: border-box;
}

.own-bubble .reply-container {
	background-color: rgba(255, 255, 255, 0.2);
	border-left-color: white;
}

.reply-header {
	font-size: 12px;
	margin-bottom: 2px;
}

.reply-content {
	font-size: 14px;
	opacity: 0.8;
}

.reply-with-photo {
	display: flex;
	align-items: center;
	gap: 8px;
}

.reply-photo {
	width: 40px;
	height: 40px;
	border-radius: 6px;
	object-fit: cover;
	flex-shrink: 0;
}

.reply-text-content {
	flex: 1;
	min-width: 0; /* Allows text to shrink and wrap if needed */
}

.message-content {
	display: flex;
	flex-direction: column;
	width: 100%; /* Ensure content takes full bubble width */
}

.message-photo {
	max-width: 300px; /* Increase photo max width */
	max-height: 300px;
	min-width: 150px; /* Add minimum width for photos */
	border-radius: 12px;
	margin-bottom: 8px;
	object-fit: cover;
}

.message-text {
	margin: 0;
	font-size: 16px;
	line-height: 1.4; /* Slightly increase line height */
	word-wrap: break-word;
	white-space: pre-wrap; /* Preserve whitespace and line breaks */
}

.message-metadata {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-top: 4px;
	font-size: 11px;
	opacity: 0.7;
	width: 100%; /* Ensure metadata spans full width */
}

.timestamp {
	margin-right: 8px;
}

.attachment-indicator {
	font-style: italic;
}

.forwarded-indicator {
  display: flex;
  align-items: center;
  font-size: 12px;
  font-style: italic;
  opacity: 0.8;
  margin-bottom: 8px;
  color: #666;
}

.own-bubble .forwarded-indicator {
  color: rgba(255, 255, 255, 0.8);
}

.forwarded-icon {
  margin-right: 4px;
  font-size: 14px;
}

.forwarded-text {
  font-weight: 500;
}
</style>
