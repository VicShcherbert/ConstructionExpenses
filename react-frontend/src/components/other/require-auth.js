import { Spinner } from "@chakra-ui/react";
import { Navigate, Outlet } from "react-router-dom";

export const RequireAuth = ({ isAuthed, authLoading }) => {
  if (authLoading) return <Spinner />;
  return isAuthed ? <Outlet /> : <Navigate to="/projects" replace />;
}
