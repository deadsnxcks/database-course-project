import type { ReactNode } from "react";
import { useNavigate } from "react-router-dom";

type Props = {
  title: string;
  children: ReactNode;
};

export default function Layout({ title, children }: Props) {
  const navigate = useNavigate();

  return (
    <div style={{ padding: "20px" }}>
      <button onClick={() => navigate("/")}>← Назад</button>
      <h1>{title}</h1>
      <hr />
      {children}
    </div>
  );
}
