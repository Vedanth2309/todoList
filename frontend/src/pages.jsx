import { useEffect, useState, useCallback } from "react";

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  LineChart,
  Line,
  PieChart,
  Pie,
  Cell,
} from "recharts";

import { Trash2, Check, Plus } from "lucide-react";

import api from "./api";

const today = () => new Date().toISOString().slice(0, 10);

const COLORS = [
  "#0f766e",
  "#b45309",
  "#be123c",
  "#4d7c0f",
  "#6d28d9",
  "#0369a1",
];

const tagsOf = (s) =>
  (s || "")
    .split(",")
    .map((t) => t.trim())
    .filter(Boolean);

function useList(path, params) {
  const [items, set] = useState(null);
  const [err, setErr] = useState("");

  const load = useCallback(
    () =>
      api
        .get(path, { params })
        .then(set)
        .catch((e) => setErr(e.message)),
    [path, JSON.stringify(params)],
  );

  useEffect(() => {
    load();
  }, [load]);

  return { items, err, load };
}

const State = ({ items, err, empty }) =>
  err ? (
    <p className="text-red-600">{err}</p>
  ) : !items ? (
    <p className="text-stone-500">Loading…</p>
  ) : items.length === 0 ? (
    <p className="text-stone-500">{empty}</p>
  ) : null;

export function Dashboard() {
  const [d, setD] = useState(null);
  const [err, setErr] = useState("");

  useEffect(() => {
    api
      .get("/analytics/dashboard")
      .then(setD)
      .catch((e) => setErr(e.message));
  }, []);

  if (err) return <p className="text-red-600">{err}</p>;
  if (!d) return <p className="text-stone-500">Loading…</p>;

  const stats = [
    ["Tasks done", d.completedTasks],
    ["Completion", Math.round(d.completionPct) + "%"],
    ["Due today", d.dueToday],
    ["Overdue", d.overdue],
    ["Habits today", `${d.habitsDoneToday}/${d.habits}`],
    ["Diary entries", d.diaryEntries],
  ];

  const pie = (data, title) => (
    <div className="card">
      <h3 className="font-semibold mb-2">{title}</h3>

      <ResponsiveContainer height={200}>
        <PieChart>
          <Pie
            data={data.map((x) => ({
              name: x._id || "None",
              value: x.count,
            }))}
            dataKey="value"
            nameKey="name"
            outerRadius={70}
            label
          >
            {data.map((_, i) => (
              <Cell key={i} fill={COLORS[i % 6]} />
            ))}
          </Pie>

          <Tooltip />
        </PieChart>
      </ResponsiveContainer>
    </div>
  );

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Today</h1>

      <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
        {stats.map(([l, v]) => (
          <div key={l} className="card">
            <div className="text-2xl font-bold">{v}</div>
            <div className="text-sm text-stone-500">{l}</div>
          </div>
        ))}
      </div>

      <div className="card">
        <h3 className="font-semibold mb-2">Tasks completed, last 14 days</h3>

        <ResponsiveContainer height={220}>
          <LineChart
            data={d.completedOverTime.map((x) => ({
              day: x._id.slice(5),
              n: x.count,
            }))}
          >
            <XAxis dataKey="day" />
            <YAxis allowDecimals={false} />
            <Tooltip />
            <Line dataKey="n" stroke="#0f766e" strokeWidth={2} />
          </LineChart>
        </ResponsiveContainer>
      </div>

      <div className="grid md:grid-cols-3 gap-3">
        {pie(d.byPriority, "By priority")}
        {pie(d.byCategory, "By category")}
        {pie(d.moods, "Mood")}
      </div>
    </div>
  );
}

export function Tasks() {
  const [status, setStatus] = useState("");
  const { items, err, load } = useList("/tasks", { status });

  const [f, setF] = useState({
    title: "",
    priority: "MEDIUM",
    dueDate: "",
    category: "",
    tags: "",
  });

  const add = async (e) => {
    e.preventDefault();

    if (!f.title) return;

    await api.post("/tasks", {
      ...f,
      tags: tagsOf(f.tags),
      status: "TODO",
    });

    setF({ ...f, title: "" });
    load();
  };

  const toggle = async (t) => {
    try {
      const response = await api.put(`/tasks/${t._id}`, {
        status: t.status === "COMPLETED" ? "TODO" : "COMPLETED",
      });

      console.log("UPDATE RESPONSE:", response);

      await load();
    } catch (error) {
      console.error("UPDATE FAILED:", error);
      console.error("SERVER RESPONSE:", error.response?.data);

      alert(
        error.response?.data?.message ||
          error.message ||
          "Failed to update task",
      );
    }
  };
  const del = (t) =>
    confirm("Delete this task?") && api.delete("/tasks/" + t._id).then(load);

  const pc = {
    HIGH: "bg-red-100 text-red-800",
    MEDIUM: "bg-amber-100 text-amber-800",
    LOW: "bg-stone-200 text-stone-700",
  };

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Tasks</h1>

      <form onSubmit={add} className="card grid md:grid-cols-6 gap-2">
        <input
          className="inp md:col-span-2"
          placeholder="New task"
          value={f.title}
          onChange={(e) => setF({ ...f, title: e.target.value })}
        />

        <select
          className="inp"
          value={f.priority}
          onChange={(e) => setF({ ...f, priority: e.target.value })}
        >
          <option>LOW</option>
          <option>MEDIUM</option>
          <option>HIGH</option>
        </select>

        <input
          className="inp"
          type="date"
          value={f.dueDate}
          onChange={(e) => setF({ ...f, dueDate: e.target.value })}
        />

        <input
          className="inp"
          placeholder="tags, comma"
          value={f.tags}
          onChange={(e) => setF({ ...f, tags: e.target.value })}
        />

        <button
          type="submit"
          className="btn flex items-center justify-center gap-1"
        >
          <Plus size={14} />
          Add task
        </button>
      </form>

      <select
        className="inp w-auto"
        value={status}
        onChange={(e) => setStatus(e.target.value)}
      >
        <option value="">All</option>
        <option>TODO</option>
        <option>IN_PROGRESS</option>
        <option>COMPLETED</option>
      </select>

      <State
        items={items}
        err={err}
        empty="No tasks yet. Add your first one above."
      />

      {items?.map((t) => (
        <div key={t._id} className="card flex items-center gap-3">
          <button
            type="button"
            onClick={() => toggle(t)}
            className={`size-6 rounded-full border transition-all hover:size-7 grid place-items-center ${
              t.status === "COMPLETED"
                ? "bg-teal-700 border-teal-700 text-white"
                : ""
            }`}
          >
            {t.status === "COMPLETED" && <Check size={14} />}
          </button>

          <div className="flex-1">
            <div
              className={
                t.status === "COMPLETED" ? "line-through text-stone-400" : ""
              }
            >
              {t.title}
            </div>

            <div className="text-xs text-stone-500">
              {t.dueDate} {t.tags?.map((g) => "#" + g).join(" ")}
            </div>
          </div>

          <span className={`text-xs rounded px-2 py-0.5 ${pc[t.priority]}`}>
            {t.priority}
          </span>

          <button type="button" onClick={() => del(t)}>
            <Trash2 size={16} />
          </button>
        </div>
      ))}
    </div>
  );
}

export function Habits() {
  const { items, err, load } = useList("/habits");

  const [name, setName] = useState("");
  const [done, setDone] = useState({});

  useEffect(() => {
    items?.forEach((h) =>
      api.get(`/habits/${h._id}/checkins`).then((c) =>
        setDone((d) => ({
          ...d,
          [h._id]: c.filter((x) => x.completed).map((x) => x.date),
        })),
      ),
    );
  }, [items]);

  const days = [...Array(7)].map((_, i) =>
    new Date(Date.now() - (6 - i) * 864e5).toISOString().slice(0, 10),
  );

  const add = async (e) => {
    e.preventDefault();

    if (name) {
      await api.post("/habits", {
        name,
        frequency: "daily",
      });

      setName("");
      load();
    }
  };

  const tick = async (h, date) => {
    const on = done[h._id]?.includes(date);

    await api.post(`/habits/${h._id}/checkin`, {
      date,
      completed: !on,
    });

    setDone((d) => ({
      ...d,
      [h._id]: on
        ? d[h._id].filter((x) => x !== date)
        : [...(d[h._id] || []), date],
    }));
  };

  const streak = (h) => {
    let n = 0;

    for (let i = 0; ; i++) {
      const d = new Date(Date.now() - i * 864e5).toISOString().slice(0, 10);

      if (done[h._id]?.includes(d)) n++;
      else if (i > 0) break;
      else if (i === 0) continue;
    }

    return n;
  };

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Habits</h1>

      <form onSubmit={add} className="flex gap-2">
        <input
          className="inp"
          placeholder="New habit, e.g. Read 20 minutes"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />

        <button type="submit" className="btn">
          Add habit
        </button>
      </form>

      <State
        items={items}
        err={err}
        empty="No habits yet. Start with one small daily habit."
      />

      {items?.map((h) => (
        <div key={h._id} className="card">
          <div className="flex justify-between mb-2">
            <b>{h.name}</b>

            <span className="text-sm text-stone-500">
              {streak(h)}-day streak
            </span>
          </div>

          <div className="flex gap-2">
            {days.map((d) => (
              <button
                type="button"
                key={d}
                onClick={() => tick(h, d)}
                className={`flex-1 rounded-lg py-2 text-xs border ${
                  done[h._id]?.includes(d)
                    ? "bg-teal-700 text-white border-teal-700"
                    : ""
                }`}
              >
                {new Date(d).toLocaleDateString("en", {
                  weekday: "short",
                })}
              </button>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}

export function Collection({ kind, fields = [] }) {
  const { items, err, load } = useList("/" + kind);

  const [f, setF] = useState({
    title: "",
    content: "",
    tags: "",
    mood: "Neutral",
  });

  const add = async (e) => {
    e.preventDefault();

    if (!f.title) return;

    await api.post("/" + kind, {
      ...f,
      tags: tagsOf(f.tags),
      date: today(),
    });

    setF({
      ...f,
      title: "",
      content: "",
    });

    load();
  };

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold capitalize">{kind}</h1>

      <form onSubmit={add} className="card space-y-2">
        <input
          className="inp"
          placeholder="Title"
          value={f.title}
          onChange={(e) => setF({ ...f, title: e.target.value })}
        />

        <textarea
          className="inp"
          rows={5}
          placeholder="Write here…"
          value={f.content}
          onChange={(e) => setF({ ...f, content: e.target.value })}
        />

        <div className="flex gap-2">
          <input
            className="inp"
            placeholder="tags, comma"
            value={f.tags}
            onChange={(e) => setF({ ...f, tags: e.target.value })}
          />

          {fields.includes("mood") && (
            <select
              className="inp w-auto"
              value={f.mood}
              onChange={(e) =>
                setF({
                  ...f,
                  mood: e.target.value,
                })
              }
            >
              {[
                "Happy",
                "Calm",
                "Neutral",
                "Excited",
                "Sad",
                "Angry",
                "Stressed",
              ].map((m) => (
                <option key={m}>{m}</option>
              ))}
            </select>
          )}

          <button type="submit" className="btn">
            Save
          </button>
        </div>
      </form>

      <State
        items={items}
        err={err}
        empty="Nothing here yet. Write your first entry above."
      />

      {items?.map((n) => (
        <div key={n._id} className="card">
          <div className="flex justify-between">
            <b>{n.title}</b>

            <button
              type="button"
              onClick={() =>
                confirm("Delete?") && api.delete(`/${kind}/${n._id}`).then(load)
              }
            >
              <Trash2 size={16} />
            </button>
          </div>

          <p className="text-sm whitespace-pre-wrap mt-1">{n.content}</p>

          <div className="text-xs text-stone-500 mt-2">
            {n.date} {n.mood} {n.tags?.map((g) => "#" + g).join(" ")}
          </div>
        </div>
      ))}
    </div>
  );
}

export function SearchPage() {
  const [q, setQ] = useState("");
  const [type, setType] = useState("");
  const [r, setR] = useState(null);
  const [err, setErr] = useState("");

  useEffect(() => {
    if (!q) return setR(null);

    const t = setTimeout(
      () =>
        api
          .get("/search", {
            params: { q, type },
          })
          .then(setR)
          .catch((e) => setErr(e.message)),
      300,
    );

    return () => clearTimeout(t);
  }, [q, type]);

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Search</h1>

      <div className="flex gap-2">
        <input
          className="inp"
          placeholder="Search tasks, diary and notes"
          value={q}
          onChange={(e) => setQ(e.target.value)}
        />

        <select
          className="inp w-auto"
          value={type}
          onChange={(e) => setType(e.target.value)}
        >
          <option value="">All</option>
          <option value="task">Tasks</option>
          <option value="diary">Diary</option>
          <option value="note">Notes</option>
        </select>
      </div>

      {err && <p className="text-red-600">{err}</p>}

      {r?.length === 0 && (
        <p className="text-stone-500">No matches for “{q}”.</p>
      )}

      {r?.map((x) => (
        <div key={x.id} className="card">
          <span className="text-xs rounded bg-teal-100 text-teal-800 px-2 py-0.5 mr-2">
            {x.type}
          </span>

          <b>{x.title}</b>

          <p className="text-sm text-stone-500 truncate">{x.content}</p>
        </div>
      ))}
    </div>
  );
}
