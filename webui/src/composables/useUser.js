export const useUser = () => {
  const getUserId = () => {
    try {
      const loggedInUser = JSON.parse(localStorage.getItem('loggedInUser'));
      return loggedInUser?.userId || null;
    } catch (error) {
      console.error('Error retrieving user ID:', error);
      return null;
    }
  };

  // Keep getUsername as a function but make it always read fresh data
  const getUsername = () => {
    try {
      const loggedInUser = JSON.parse(localStorage.getItem('loggedInUser'));
      console.log('Retrieved loggedInUser:', loggedInUser); // Debugging line
      return loggedInUser?.name || null;

    } catch (error) {
      console.error('Error retrieving username:', error);
      return null;
    }
  };

  // Add a function to update user data in localStorage
  const updateUserData = (userData) => {
    try {
      const currentUser = JSON.parse(localStorage.getItem('loggedInUser') || '{}');
      const updatedUser = { ...currentUser, ...userData };
      localStorage.setItem('loggedInUser', JSON.stringify(updatedUser));
    } catch (error) {
      console.error('Error updating user data:', error);
    }
  };

  return { getUserId, getUsername, updateUserData };
};