import { useState, useEffect } from "react";
import {
  Routes,
  Route,
  NavLink,
  Navigate,
  useNavigate,
} from "react-router-dom";
import {
  LayoutDashboard,
  CheckSquare,
  Repeat,
  BookOpen,
  StickyNote,
  Search,
  Moon,
  Sun,
  LogOut,
} from "lucide-react";
import api from "./api";
import { Dashboard, Tasks, Habits, Collection, SearchPage } from "./pages";

function Auth({ mode }) {
  const nav = useNavigate();
  const [f, setF] = useState({});
  const [err, setErr] = useState("");
  const submit = async (e) => {
    e.preventDefault();
    setErr("");
    if (mode === "register" && f.password !== f.confirm)
      return setErr("Passwords do not match");
    try {
      const r = await api.post("/auth/" + mode, f);
      localStorage.setItem("token", r.token);
      nav("/");
    } catch (x) {
      setErr(x.message);
    }
  };
  const set = (k) => (e) => setF({ ...f, [k]: e.target.value });
  return (
    <div className="min-h-screen grid place-items-center p-4">
      <form onSubmit={submit} className="card w-full max-w-sm space-y-3">
        <h1 className="text-xl font-bold">
          {mode === "login" ? "Welcome back" : "Create your account"}
        </h1>
        {mode === "register" && (
          <input className="inp" placeholder="Name" onChange={set("name")} />
        )}
        <input
          className="inp"
          type="email"
          placeholder="Email"
          onChange={set("email")}
        />
        <input
          className="inp"
          type="password"
          placeholder="Password"
          onChange={set("password")}
        />
        {mode === "register" && (
          <input
            className="inp"
            type="password"
            placeholder="Confirm password"
            onChange={set("confirm")}
          />
        )}
        {err && <p className="text-sm text-red-600">{err}</p>}
        <button className="btn w-full">
          {mode === "login" ? "Log in" : "Sign up"}
        </button>
        <NavLink
          className="block text-sm text-teal-700"
          to={mode === "login" ? "/register" : "/login"}
        >
          {mode === "login"
            ? "Need an account? Sign up"
            : "Have an account? Log in"}
        </NavLink>
      </form>
    </div>
  );
}

const links = [
  ["/", "Dashboard", LayoutDashboard],
  ["/tasks", "Tasks", CheckSquare],
  ["/habits", "Habits", Repeat],
  ["/diary", "Diary", BookOpen],
  ["/notes", "Notes", StickyNote],
  ["/search", "Search", Search],
];

function Shell() {
  const [dark, setDark] = useState(
    document.documentElement.classList.contains("dark"),
  );
  useEffect(() => {
    document.documentElement.classList.toggle("dark", dark);
    localStorage.setItem("dark", dark ? "1" : "0");
  }, [dark]);
  if (!localStorage.getItem("token")) return <Navigate to="/login" />;
  const cls = ({ isActive }) =>
    `flex items-center gap-2 rounded-lg px-3 py-2 text-sm ${isActive ? "bg-teal-700 text-white" : "hover:bg-stone-200 dark:hover:bg-stone-800"}`;
  return (
    <div className="md:flex min-h-screen">
      <aside className="hidden md:flex w-56 flex-col gap-1 p-3 border-r border-stone-200 dark:border-stone-800">
        {links.map(([to, l, I]) => (
          <NavLink key={to} to={to} end={to === "/"} className={cls}>
            <I size={16} />
            {l}
          </NavLink>
        ))}
        <div className="mt-auto flex gap-2">
          <button onClick={() => setDark(!dark)} className="p-2">
            {dark ? <Sun size={16} /> : <Moon size={16} />}
          </button>
          <button
            onClick={() => {
              localStorage.removeItem("token");
              location.href = "/login";
            }}
            className="p-2"
          >
            <LogOut size={16} />
          </button>
        </div>
      </aside>
      <main className="flex-1 p-4 md:p-8 pb-20 md:pb-8 max-w-5xl">
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/tasks" element={<Tasks />} />
          <Route path="/habits" element={<Habits />} />
          <Route
            path="/diary"
            element={<Collection kind="diary" fields={["mood"]} />}
          />
          <Route path="/notes" element={<Collection kind="notes" />} />
          <Route path="/search" element={<SearchPage />} />
        </Routes>
      </main>
      <nav className="md:hidden fixed bottom-0 inset-x-0 flex justify-around border-t border-stone-200 dark:border-stone-800 bg-white dark:bg-stone-900 py-2">
        {links.map(([to, l, I]) => (
          <NavLink
            key={to}
            to={to}
            end={to === "/"}
            className={({ isActive }) => (isActive ? "text-teal-700" : "")}
          >
            <I size={20} />
          </NavLink>
        ))}
      </nav>
    </div>
  );
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Auth mode="login" />} />
      <Route path="/register" element={<Auth mode="register" />} />
      <Route path="/*" element={<Shell />} />
    </Routes>
  );
}
