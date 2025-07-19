const Failed = () => {
  return (
    <div className="flex justify-center items-center py-4">
      <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded flex items-center gap-2">
        <span className="text-xl">❌</span>
        <span>Failed to upload document. Please try again.</span>
      </div>
    </div>
  );
};

export default Failed;
