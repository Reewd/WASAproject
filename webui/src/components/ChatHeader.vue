<template>
	<div class="chat-header">
		<!-- Display the chat photo -->
		<img :src="headerImageUrl" :alt="headerImageAlt" class="chat-photo" />
		
		<!-- Display chat title and additional info -->
		<div class="chat-info">
			<h2>{{ headerTitle }}</h2>
			<!-- Show participant status for private chats -->
			<span
				v-if="!isGroup && otherParticipant"
				class="participant-status"
			>
			</span>
			<!-- Show participant count for group chats -->
			<span v-else-if="isGroup && participants" class="participant-count">
				{{ participants.length }} participants
			</span>
		</div>

		<!-- Action buttons for group settings -->
		<div class="chat-actions">
			<button
				v-if="isGroup"
				@click="openGroupSettings"
				class="group-settings-button"
			>
				Group Settings
			</button>
		</div>

		<!-- Group settings modal -->
		<GroupSettings
			v-if="showGroupSettings"
			:chat="props.chat"
			:isGroup="isGroup"
			@close="closeGroupSettings"
			@updated="handleGroupSettingsUpdated"
			@left="handleLeftGroup"
		/>
	</div>
</template>

<script setup>
import { ref, computed } from "vue";
import { useAuth } from "../composables/useAuth.js";
import { useImageUrl } from "../composables/useImageUrl.js";
import GroupSettings from "../modals/GroupSettings.vue";
import groupDefaultIcon from "/assets/icons/group-default.png";
import userDefaultIcon from "/assets/icons/user-default.png";

const { getCurrentUserId } = useAuth();
const { getImageUrl } = useImageUrl();

const props = defineProps({
	chat: {
		type: Object,
		required: false,
		default: null,
	},
	conversationPreview: {
		type: Object,
		required: false,
		default: null,
	},
});

const emits = defineEmits(["groupUpdated", "leftGroup"]);

const handleLeftGroup = () => {
	emits("leftGroup");
};

const showGroupSettings = ref(false);

const isGroup = computed(() => {
	return props.chat?.isGroup || false;
});

const participants = computed(() => {
	return props.chat?.participants || [];
});

const otherParticipant = computed(() => {
	if (isGroup.value || !participants.value.length) return null;

	return (
		participants.value.find((p) => p.userId !== getCurrentUserId()) || null
	);
});

const headerTitle = computed(() => {
	if (isGroup.value) {
		return (
			props.chat?.name || props.conversationPreview?.name || "Group Chat"
		);
	} else {
		return (
			otherParticipant.value?.username ||
			props.conversationPreview?.name ||
			"Private Chat"
		);
	}
});

const headerImageUrl = computed(() => {
	if (isGroup.value) {
		if (props.chat?.photo?.path) {
			return getImageUrl(props.chat.photo.path);
		}
		return groupDefaultIcon;
	} else {
		if (otherParticipant.value?.photo?.path) {
			return getImageUrl(otherParticipant.value.photo.path);
		}
		return userDefaultIcon;
	}
});

const headerImageAlt = computed(() => {
	if (isGroup.value) {
		return "Group Photo";
	} else {
		return `${
			otherParticipant.value?.username || "User"
		}'s Profile Picture`;
	}
});

const openGroupSettings = () => {
	showGroupSettings.value = true;
};

const closeGroupSettings = () => {
	showGroupSettings.value = false;
};

const handleGroupSettingsUpdated = () => {
	showGroupSettings.value = false;
	emits("groupUpdated");
};
</script>

<style scoped>
.chat-header {
	display: flex;
	align-items: center;
	padding: 16px 20px;
	border-bottom: 1px solid #eee;
	background-color: #f8f9fa;
	min-height: 70px;
}

.chat-photo {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	margin-right: 12px;
	object-fit: cover;
	flex-shrink: 0;
}

.chat-info {
	flex: 1;
	display: flex;
	flex-direction: column;
	min-width: 0;
}

.chat-info h2 {
	margin: 0;
	font-size: 18px;
	font-weight: 600;
	color: #333;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.participant-status,
.participant-count {
	font-size: 12px;
	color: #666;
	margin-top: 2px;
}

.chat-actions {
	display: flex;
	align-items: center;
	gap: 8px;
}

.group-settings-button {
	background-color: #007bff;
	color: white;
	border: none;
	border-radius: 6px;
	padding: 8px 16px;
	font-size: 14px;
	font-weight: 500;
	cursor: pointer;
	transition: background-color 0.2s;
	white-space: nowrap;
}

.group-settings-button:hover {
	background-color: #0056b3;
}

.group-settings-button:active {
	background-color: #004085;
}
</style>
