const Section = ({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) => {
  return (
    <div className="space-y-2">
      <div className="font-medium">{title}</div>
      {children}
    </div>
  );
};

export default Section;
