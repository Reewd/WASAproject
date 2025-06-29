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

  const getUsername = () => {
    try {
      const loggedInUser = JSON.parse(localStorage.getItem('loggedInUser'));
      return loggedInUser?.name || null;
    } catch (error) {
      console.error('Error retrieving username:', error);
      return null;
    }
  };

  return { getUserId, getUsername };
};