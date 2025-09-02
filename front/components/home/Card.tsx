import { ReactNode } from "react";

const Card = ({
  children,
  className,
}: {
  children: ReactNode;
  className?: string;
}) => {
  return (
    <div className={`${className}`}>
      <div className="bg-white/80 backdrop-blur rounded-2xl shadow border border-purple-100 overflow-hidden">
        {children}
      </div>
    </div>
  );
};

export default Card;
