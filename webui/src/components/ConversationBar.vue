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
import { useUser } from "../composables/useUser.js"; // Import the composable
import NewConversation from "../modals/NewConversation.vue"; // Import the modal
import newConversationIcon from "/assets/icons/new-conversation.png";

const conversations = ref([]);
const sidebar = ref(null);
let isResizing = false;
const showModal = ref(false); // State to control modal visibility
const emits = defineEmits(["selectConversation"]); // Emit selected conversation
const pollingInterval = ref(null);
const POLLING_DELAY = 5000; // Poll every 5 seconds

const { getUserId } = useUser(); // Use the composable to retrieve the userId

const fetchConversations = async () => {
	try {
		const userId = getUserId(); // Retrieve the userId using the composable

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
	width: 5px; /* Set a visible width */
	height: 100%;
	position: absolute;
	top: 0;
	right: 0;
	background-color: transparent; /* Make it invisible but still clickable */
	cursor: ew-resize; /* Show resize cursor */
}

.resize-handle:hover {
	background-color: #ccc; /* Highlight the handle on hover */
}

.new-conversation-button {
	position: absolute;
	top: 10px;
	right: 10px;
	background: none; /* Removes background color */
	border: none; /* Removes border */
	cursor: pointer;
	padding: 0; /* Removes padding */
	width: 40px; /* Set width for circular shape */
	height: 40px; /* Set height for circular shape */
	border-radius: 50%; /* Makes the button circular */
	display: flex; /* Centers the image */
	align-items: center; /* Centers the image */
	justify-content: center; /* Centers the image */
	transition: box-shadow 0.3s ease; /* Smooth shadow transition */
}

.new-conversation-button:hover {
	box-shadow: 0 0 10px rgba(0, 0, 0, 0.3); /* Adds circular shadow on hover */
	opacity: 0.8; /* Slightly darkens the button on hover */
}

.new-conversation-button img {
	width: 24px; /* Adjust image size */
	height: 24px; /* Adjust image size */
}

.selected {
	background-color: #333; /* Dark gray background */
	color: #fff; /* White text for better contrast */
}
</style>
