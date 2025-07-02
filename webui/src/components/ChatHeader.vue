<template>
	<div class="chat-header">
		<img :src="headerImageUrl" :alt="headerImageAlt" class="chat-photo" />
		<div class="chat-info">
			<h2>{{ headerTitle }}</h2>
			<span
				v-if="!isGroup && otherParticipant"
				class="participant-status"
			>
			</span>
			<span v-else-if="isGroup && participants" class="participant-count">
				{{ participants.length }} participants
			</span>
		</div>

		<div class="chat-actions">
			<button
				v-if="isGroup"
				@click="openGroupSettings"
				class="group-settings-button"
			>
				Group Settings
			</button>
		</div>

		<!-- Modal managed directly by ChatHeader -->
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
import GroupSettings from "../modals/GroupSettings.vue"; // Import the modal
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

const emits = defineEmits(["groupUpdated", "leftGroup"]); // Only emit when something actually changes


const handleLeftGroup = () => {
	// Emit event to parent component when user leaves group
	emits("leftGroup");
};

// Modal state
const showGroupSettings = ref(false);

// Computed properties
const isGroup = computed(() => {
	return props.chat?.isGroup || false;
});

const participants = computed(() => {
	return props.chat?.participants || [];
});

// For private conversations, get the other participant (not the current user)
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
		// Group conversation - use group photo or default group icon
		if (props.chat?.photo?.path) {
			return getImageUrl(props.chat.photo.path);
		}
		return groupDefaultIcon;
	} else {
		// Private conversation - use other participant's photo or default user icon
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
	// Emit to parent that group was updated so it can refresh data
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
	min-width: 0; /* Allow text to truncate */
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
