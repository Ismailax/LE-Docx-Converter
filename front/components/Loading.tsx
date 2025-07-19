const Loading = () => {
  return (
    <div className="flex flex-col justify-center items-center py-4">
      <div className="loader" />
      <div className="ml-2 text-blue-600">Processing document...</div>
    </div>
  );
};

export default Loading;
