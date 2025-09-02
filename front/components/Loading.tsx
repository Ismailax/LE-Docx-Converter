const Loading = () => {
  return (
    <div className="flex flex-col justify-center items-center space-y-3">
      <div className="h-12 w-12 rounded-full border-5 border-gray-200 border-t-purple-600 animate-spin" />
      <div className="ml-2 text-purple-600 text-lg">Processing document...</div>
    </div>
  );
};

export default Loading;
