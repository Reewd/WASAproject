<template>
    <div class="conversation-preview">
        <img
            :src="conversationPhotoUrl"
            alt="Conversation Photo"
            class="conversation-photo"
        />
        <div class="conversation-details">
            <h3 class="conversation-name">
                {{ displayName }}
            </h3>
            <p class="last-message" v-if="conversation.lastMessage">
                <strong>{{ conversation.lastMessage.sentBy.username }}:</strong>
                <span
                    v-if="
                        conversation.lastMessage.photo &&
                        conversation.lastMessage.text
                    "
                >
                    📷 {{ conversation.lastMessage.text }}
                </span>
                <span v-else-if="conversation.lastMessage.photo">
                    📷 Photo
                </span>
                <span v-else>
                    {{ conversation.lastMessage.text }}
                </span>
            </p>
            <p class="last-message" v-else>No messages yet</p>
            <p v-if="conversation.lastMessage">
                <span class="last-message timestamp">{{
                    formatTimestamp(conversation.lastMessage.timestamp)
                }}</span>
            </p>
        </div>
    </div>
</template>

<script setup>
import { computed } from "vue";
import { useAuth } from "../composables/useAuth.js";
import { useImageUrl } from "../composables/useImageUrl.js";
import groupDefaultIcon from "/assets/icons/group-default.png";
import userDefaultIcon from "/assets/icons/user-default.png";

const { user } = useAuth();
const { getImageUrl } = useImageUrl();

const props = defineProps({
	conversation: {
		type: Object,
		required: true,
	},
});

const otherParticipant = computed(() => {
	if (props.conversation.isGroup) return null;

	return (
		props.conversation.participants?.find(
			(participant) => participant.userId !== user.value.userId
		) || null
	);
});

const displayName = computed(() => {
	if (props.conversation.isGroup) {
		return props.conversation.name;
	} else {
		return otherParticipant.value
			? otherParticipant.value.username
			: props.conversation.name;
	}
});

const conversationPhotoUrl = computed(() => {
	if (props.conversation.isGroup) {
		if (props.conversation.photo?.path) {
			return getImageUrl(props.conversation.photo.path);
		}
		return groupDefaultIcon;
	} else {
		if (otherParticipant.value?.photo?.path) {
			return getImageUrl(otherParticipant.value.photo.path);
		}
		return userDefaultIcon;
	}
});

const formatTimestamp = (timestamp) => {
	const date = new Date(timestamp);
	return `${date.toLocaleDateString([], {
		year: "numeric",
		month: "short",
		day: "numeric",
	})} ${date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}`;
};
</script>

<style scoped>
.conversation-preview {
    display: flex;
    align-items: center;
    padding: 10px;
    border-radius: 5px;
    min-width: 0;
}

.conversation-preview:hover {
    background-color: #f0f0f0;
}

.conversation-photo {
    width: 50px;
    height: 50px;
    border-radius: 50%;
    margin-right: 10px;
    object-fit: cover;
    flex-shrink: 0;
}

.conversation-preview.selected {
    background-color: #333;
    color: #fff;
}

.conversation-details {
    flex: 1;
    min-width: 0;
    overflow: hidden;
}

.conversation-name {
    font-size: 16px;
    font-weight: bold;
    margin: 0;
}

.last-message {
    font-size: 14px;
    color: #666;
    margin: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
}

.conversation-preview.selected .last-message {
    color: #ccc;
}
</style>