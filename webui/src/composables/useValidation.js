import { ref, computed } from 'vue';

export function useValidation() {
  // Generic validation for names (usernames, group names, etc.)
  const validateName = (value, fieldName = 'Name') => {
    if (!value || value.trim() === '') {
      return `${fieldName} is required`;
    }
    if (value.length < 3 || value.length > 16) {
      return `${fieldName} must be between 3 and 16 characters`;
    }
    return '';
  };

  // Composable for username validation
  const useUsernameValidation = (initialValue = '') => {
    const username = ref(initialValue);
    const usernameError = ref('');

    const validateUsername = () => {
      usernameError.value = validateName(username.value, 'Username');
    };

    const isUsernameValid = computed(() => {
      return usernameError.value === '' && username.value.trim() !== '';
    });

    return {
      username,
      usernameError,
      validateUsername,
      isUsernameValid
    };
  };

  // Composable for group name validation
  const useGroupNameValidation = (initialValue = '') => {
    const groupName = ref(initialValue);
    const groupNameError = ref('');

    const validateGroupName = () => {
      groupNameError.value = validateName(groupName.value, 'Group name');
    };

    const isGroupNameValid = computed(() => {
      return groupNameError.value === '' && groupName.value.trim() !== '';
    });

    return {
      groupName,
      groupNameError,
      validateGroupName,
      isGroupNameValid
    };
  };

  return {
    validateName,
    useUsernameValidation,
    useGroupNameValidation
  };
}
