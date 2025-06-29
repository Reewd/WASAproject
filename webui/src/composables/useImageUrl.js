export function useImageUrl() {
  const getImageUrl = (path) => {
    if (!path) return null;
    return `${__API_URL__}/${path}`;
  };

  return {
    getImageUrl
  };
}