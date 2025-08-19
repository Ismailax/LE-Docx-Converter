const Loading = () => {
  return (
    <div className="flex flex-col justify-center items-center space-y-3">
      <div className="loader" />
      <div className="ml-2 text-purple-600 text-lg">Processing document...</div>
    </div>
  );
};

export default Loading;
