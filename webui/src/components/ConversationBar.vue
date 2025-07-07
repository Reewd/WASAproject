<template>
	<div class="resizable-sidebar" ref="sidebar">
		<br /><br />
		<!-- New Conversation Button -->
		<button class="new-conversation-button" @click="showModal = true">
			<img
				:src="newConversationIcon"
				alt="New Conversation"
			/>
		</button>

		<template v-if="conversations.length > 0">
			<div
				v-for="conversation in conversations"
				:key="conversation.conversationId"
				@click="selectConversation(conversation)"
				:class="'conversation-preview'"
			>
				<ConversationPreview :conversation="conversation" />
			</div>
		</template>
		<div class="resize-handle" @mousedown="startResizing"></div>

		<!-- New Conversation Modal -->
		<NewConversation
			v-if="showModal"
			@close="showModal = false"
			@conversationCreated="fetchConversations"
			:existingConversations="conversations"
		/>
	</div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import axios from "../services/axios.js";
import ConversationPreview from "./ConversationPreview.vue";
import { useAuth } from "../composables/useAuth.js";
import NewConversation from "../modals/NewConversation.vue"; // Import the modal
import newConversationIcon from "/assets/icons/new-conversation.png";

const conversations = ref([]);
const sidebar = ref(null);
let isResizing = false;
const showModal = ref(false); // State to control modal visibility
const emits = defineEmits(["selectConversation"]); // Emit selected conversation
const pollingInterval = ref(null);
const POLLING_DELAY = 5000; // Poll every 5 seconds

const { user } = useAuth();

const fetchConversations = async () => {
	try {
		const userId = user.value.userId; // Retrieve the userId using the composable
		console.log("Fetching conversations for user ID:", userId);

		if (!userId) {
			console.error("User ID not found");
			return;
		}

		const response = await axios.get("/conversations", {
			headers: {
				Authorization: userId,
			},
		});
		
		// Check if there are changes before updating to avoid unnecessary re-renders
		if (JSON.stringify(conversations.value) !== JSON.stringify(response.data.conversations)) {
			conversations.value = response.data.conversations;
			console.log("Conversations updated:", conversations.value);
		}
	} catch (error) {
		console.error("Error fetching conversations:", error);
	}
};

// Start polling for conversation updates
const startPolling = () => {
	// Clear any existing interval first
	stopPolling();
	
	// Start new polling interval
	pollingInterval.value = setInterval(() => {
		fetchConversations();
	}, POLLING_DELAY);
};

// Stop polling
const stopPolling = () => {
	if (pollingInterval.value) {
		clearInterval(pollingInterval.value);
		pollingInterval.value = null;
	}
};

const selectConversation = (conversation) => {
	emits("selectConversation", conversation); // Emit the selected conversation
};

// Resizing logic
const startResizing = () => {
	isResizing = true;
	document.addEventListener("mousemove", resizeSidebar);
	document.addEventListener("mouseup", stopResizing);
};

const resizeSidebar = (event) => {
	if (isResizing && sidebar.value) {
		const sidebarElement = sidebar.value;
		const newWidth = Math.max(
			200,
			Math.min(
				500,
				event.clientX - sidebarElement.getBoundingClientRect().left
			)
		);
		sidebarElement.style.width = `${newWidth}px`;
	}
};

const stopResizing = () => {
	isResizing = false;
	document.removeEventListener("mousemove", resizeSidebar);
	document.removeEventListener("mouseup", stopResizing);
};

// Lifecycle hooks
onMounted(() => {
	fetchConversations();
	startPolling();
	document.addEventListener("mouseup", stopResizing);
});

onUnmounted(() => {
	stopPolling();
	document.removeEventListener("mousemove", resizeSidebar);
	document.removeEventListener("mouseup", stopResizing);
});
defineExpose({
	fetchConversations,
	startPolling,
	stopPolling
});
</script>

<style scoped>
.resizable-sidebar {
	width: 300px;
	min-width: 200px;
	max-width: 500px;
	overflow: auto;
	border: 1px solid #ccc;
	padding: 10px;
	position: relative;
}

.conversation-preview {
	margin-bottom: 5px;
	padding-bottom: 5px;
}

.resize-handle {
	width: 5px;
	height: 100%;
	position: absolute;
	top: 0;
	right: 0;
	background-color: transparent;
	cursor: ew-resize;
}

.resize-handle:hover {
	background-color: #ccc;
}

.new-conversation-button {
	position: absolute;
	top: 10px;
	right: 10px;
	background: none;
	border: none;
	cursor: pointer;
	padding: 0;
	width: 40px;
	height: 40px;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	transition: box-shadow 0.3s ease;
}

.new-conversation-button:hover {
	box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
	opacity: 0.8;
}

.new-conversation-button img {
	width: 24px;
	height: 24px;
}

.selected {
	background-color: #333;
	color: #fff;
}
</style>
