(() => {
  "use strict";
  const KEY = "bakery-shared-table-v1";
  const $ = (selector) => document.querySelector(selector);
  const esc = (value) =>
    String(value ?? "").replace(
      /[&<>"']/g,
      (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" })[c],
    );
  const paths = {
    notes: "M5 3h14v18H5z M8 8h8 M8 12h8 M8 16h5",
    review: "M7 3v12a4 4 0 0 0 4 4h6 M17 5v14 M14 8l3-3 3 3 M4 5h6",
    branch: "M6 4v16 M6 14c8 0 12-3 12-9 M3 4a3 3 0 1 0 6 0 M3 20a3 3 0 1 0 6 0 M15 5a3 3 0 1 0 6 0",
    activity: "M3 12h4l3-7 4 14 3-7h4",
    team: "M9 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8 M2 21v-2a7 7 0 0 1 14 0v2 M17 5a4 4 0 0 1 0 7 M19 16a5 5 0 0 1 3 5",
    terminal: "M3 5h18v15H3z M7 9l3 3-3 3 M13 15h4",
    book: "M3 4c3-1 6 0 9 2 3-2 6-3 9-2v15c-3-1-6 0-9 2-3-2-6-3-9-2z M12 6v15",
    grid: "M3 3h7v7H3z M14 3h7v7h-7z M3 14h7v7H3z M14 14h7v7h-7z",
    network:
      "M12 12 5 5 M12 12l7-7 M12 12 5 20 M12 12l14 7 M2 2h6v6H2z M16 2h6v6h-6z M2 17h6v6H2z M16 16h6v6h-6z",
    pin: "M8 3h8l-1 7 4 4H5l4-4z M12 14v8",
    link: "M10 8l3-3a4 4 0 0 1 6 6l-3 3 M14 16l-3 3a4 4 0 0 1-6-6l3-3 M8 16l8-8",
    arrow: "M4 12h16 M14 6l6 6-6 6",
    check: "m5 12 4 4L19 6",
    clock: "M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18 M12 7v6l4 2",
    folder: "M3 6h7l2 3h9v12H3z",
    download: "M12 3v12 M7 10l5 5 5-5 M3 17v4h18v-4",
    spark: "m12 3 2 6 7 3-7 2-2 7-2-7-7-2 7-3z",
    shield: "M12 3l8 3v6c0 5-8 9-8 9s-8-4-8-9V6z M8 12l3 3 5-6",
    edit: "m4 16 12-12 4 4L8 20l-5 1z M13 7l4 4",
  };
  function icon(name) {
    return (
      '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="' +
      (paths[name] || paths.notes) +
      '"/></svg>'
    );
  }
  function seed() {
    return {
      schema: 1,
      view: "memory",
      project: "all",
      kind: "All",
      layout: "cards",
      branch: "main",
      search: "",
      focused: "",
      reviewFilter: "open",
      selectedReview: "",
      selectedBranch: "",
      pins: ["codes"],
      notes: [
        {
          id: "codes",
          title: "Account codes must remain strings.",
          short: "Keep the leading zeros",
          kind: "Decision",
          project: "Payments",
          body: "007 and 7 are different account identifiers. Preserve leading zeros and letter case when exporting account codes. Normalize outer whitespace only.",
          why: "A receiver rejected identifiers after the exporter converted them to numbers.",
          source: "tests/test_account_export.py",
          commit: "a8c41d2",
          verified: "2 Nov",
          revision: 2,
          author: "Priya Shah",
          initials: "PS",
          links: ["tenant", "cart"],
          evidence: 'assert export_code(" 007 ") == "007"\nassert export_code("aB09") == "aB09"',
        },
        {
          id: "retries",
          title: "Retry payments only when acceptance is known.",
          short: "A safer second try",
          kind: "Decision",
          project: "Payments",
          body: "Retry explicit HTTP 429 and 503 responses, using the original idempotency key. Treat other outcomes as a separate investigation.",
          why: "Retrying a charge that was already accepted can create a duplicate payment.",
          source: "tests/billing/retries.test.ts",
          commit: "98a2c11",
          verified: "1 Nov",
          revision: 2,
          author: "Alex Morgan",
          initials: "AM",
          links: ["webhooks", "sandbox"],
          evidence: "expect(retryable(429)).toBe(true)\nexpect(retryable(400)).toBe(false)",
        },
        {
          id: "webhooks",
          title: "Deduplicate webhook events before processing.",
          short: "One event, one credit",
          kind: "Decision",
          project: "Payments",
          body: "Store the provider event ID before applying a credit. Repeated delivery of the same event must not credit the account again.",
          why: "Webhook delivery is at least once. A successful delivery may still be repeated.",
          source: "tests/billing/webhooks.test.ts",
          commit: "f09de44",
          verified: "3 Nov",
          revision: 1,
          author: "Sam Kim",
          initials: "SK",
          links: ["retries", "tenant"],
          evidence: "process(event)\nprocess(event)\nexpect(credits.count()).toBe(1)",
        },
        {
          id: "tenant",
          title: "Billing settings belong to the tenant.",
          short: "A tenant has its own tab",
          kind: "Knowledge",
          project: "Payments",
          body: "Resolve billing settings through the tenant, not through the currently signed-in user. A user can belong to more than one tenant.",
          why: "Account-level settings leaked between tenants when a user switched organizations.",
          source: "src/tenants/billing.ts",
          commit: "abc751c",
          verified: "30 Oct",
          revision: 1,
          author: "Priya Shah",
          initials: "PS",
          links: ["codes", "webhooks"],
          evidence: "billingSettings = tenant.billingSettings\n// The current user is not the billing owner.",
        },
        {
          id: "workers",
          title: "Restart queue workers after every deployment.",
          short: "Fresh code needs fresh workers",
          kind: "Lesson",
          project: "Platform",
          body: "Long-running queue workers keep old application code in memory. Restart the worker pool after the new deployment is healthy.",
          why: "A background job used yesterday’s code after the application deployed successfully.",
          source: ".github/workflows/deploy.yml",
          commit: "7ba920e",
          verified: "20 Oct",
          revision: 2,
          author: "Jules Chen",
          initials: "JC",
          links: ["sandbox"],
          stale: true,
          evidence:
            "deploy → wait for health check → restart workers\n// Referenced workflow changed after verification.",
        },
        {
          id: "sandbox",
          title: "Use the provider sandbox for integration tests.",
          short: "Practice on the small oven",
          kind: "Lesson",
          project: "Payments",
          body: "Use the provider sandbox for integration tests that verify retry signatures. The static fixture misses the receiver’s retry behavior.",
          why: "An integration passed with the fixture but failed against the real sandbox.",
          source: "tests/integration/provider.test.ts",
          commit: "a8c41d2",
          verified: "4 Nov",
          revision: 1,
          author: "Sam Kim",
          initials: "SK",
          links: ["retries", "workers"],
          evidence: "provider = createSandboxClient()\nexpect(await provider.verifyRetry()).toBe(true)",
        },
        {
          id: "cart",
          title: "Cart prices are confirmed at checkout.",
          short: "Check the price before the bite",
          kind: "Decision",
          project: "Storefront",
          body: "Revalidate cart prices when checkout begins. The cart is a preview, not a promise that an expired offer is still valid.",
          why: "An offer expired while a cart was left open overnight.",
          source: "tests/cart/checkout.test.ts",
          commit: "d18ea44",
          verified: "3 Nov",
          revision: 1,
          author: "Jules Chen",
          initials: "JC",
          links: ["money", "codes"],
          evidence: "expect(checkout(cart).price).toEqual(currentPrice)",
        },
        {
          id: "money",
          title: "Store money in integer minor units.",
          short: "A penny is a whole thing",
          kind: "Knowledge",
          project: "Storefront",
          body: "Persist monetary amounts in integer minor units, alongside their currency. Format for display only at the boundary.",
          why: "Floating-point rounding produced inconsistent line-item totals.",
          source: "src/money/amount.ts",
          commit: "cd1928a",
          verified: "2 Nov",
          revision: 1,
          author: "Alex Morgan",
          initials: "AM",
          links: ["cart"],
          evidence: '{ amount: 1299, currency: "USD" }\n// Display as 12.99 only at the UI boundary.',
        },
      ],
      branches: [
        {
          name: "billing/receipt-check",
          project: "Payments",
          status: "open",
          author: "Priya Shah",
          description: "Make ambiguous payment timeouts explicit.",
        },
        {
          name: "checkout/price-window",
          project: "Storefront",
          status: "open",
          author: "Jules Chen",
          description: "Clarify when the displayed price is revalidated.",
        },
        {
          name: "deploy/worker-pool",
          project: "Platform",
          status: "open",
          author: "Sam Kim",
          description: "Reverify the lesson against the new deployment workflow.",
        },
      ],
      proposals: [
        {
          id: "p-retries",
          noteId: "retries",
          branch: "billing/receipt-check",
          status: "pending",
          baseRevision: 2,
          baseBody:
            "Retry explicit HTTP 429 and 503 responses, using the original idempotency key. Treat other outcomes as a separate investigation.",
          title: "Retry payments only when acceptance is known.",
          body: "Retry explicit HTTP 429 and 503 responses using the original idempotency key. After a timeout, look up the provider receipt before retrying: the charge may already have been accepted.",
          source: "tests/billing/retries.test.ts",
          reason: "A timeout is an unknown outcome. This change makes the receipt check explicit.",
          author: "Priya Shah",
          agent: "Claude Code",
          kind: "Decision",
          project: "Payments",
          comments: [
            {
              author: "Priya",
              text: "Added the timeout regression case to the evidence.",
              when: "18 minutes ago",
            },
          ],
        },
        {
          id: "p-cart",
          noteId: "cart",
          branch: "checkout/price-window",
          status: "pending",
          baseRevision: 1,
          baseBody:
            "Revalidate cart prices when checkout begins. The cart is a preview, not a promise that an expired offer is still valid.",
          title: "Cart prices are confirmed at checkout.",
          body: "Preserve the displayed price while browsing the cart. At checkout, revalidate price and availability; ask the customer to confirm any changed total before payment.",
          source: "tests/cart/checkout.test.ts",
          reason: "The customer should see and confirm a changed total before being charged.",
          author: "Jules Chen",
          agent: "Codex",
          kind: "Decision",
          project: "Storefront",
          comments: [],
        },
        {
          id: "p-workers",
          noteId: "workers",
          branch: "deploy/worker-pool",
          status: "pending",
          baseRevision: 2,
          baseBody:
            "Long-running queue workers keep old application code in memory. Restart the worker pool after the new deployment is healthy.",
          title: "Restart queue workers after every deployment.",
          body: "After the new deployment passes its health check, roll the queue worker pool. Wait for in-flight jobs to finish before stopping the old workers.",
          source: ".github/workflows/deploy.yml",
          reason: "The new worker-pool deployment drains jobs instead of terminating them.",
          author: "Sam Kim",
          agent: "Codex",
          kind: "Lesson",
          project: "Platform",
          comments: [],
        },
      ],
      members: [
        { name: "Alex Morgan", email: "alex@example.invalid", role: "Owner", initials: "AM" },
        { name: "Priya Shah", email: "priya@example.invalid", role: "Reviewer", initials: "PS" },
        { name: "Sam Kim", email: "sam@example.invalid", role: "Contributor", initials: "SK" },
        { name: "Jules Chen", email: "jules@example.invalid", role: "Reviewer", initials: "JC" },
      ],
      events: [
        {
          title: "Priya proposed a safer timeout policy.",
          detail: "Payments · billing/receipt-check · via Claude Code",
          when: "18 minutes ago",
          initials: "PS",
          kind: "proposal",
        },
        {
          title: "Jules clarified the checkout price window.",
          detail: "Storefront · checkout/price-window · via Codex",
          when: "46 minutes ago",
          initials: "JC",
          kind: "proposal",
        },
        {
          title: "The deployment lesson needs another look.",
          detail: "Platform · evidence changed since verification",
          when: "2 hours ago",
          initials: "SK",
          kind: "evidence",
        },
        {
          title: "Alex merged account-code guidance.",
          detail: "Payments · revision 2 is shared on main",
          when: "Yesterday",
          initials: "AM",
          kind: "merge",
        },
      ],
    };
  }
  let s = seed();
  try {
    const stored = JSON.parse(localStorage.getItem(KEY));
    if (
      stored?.schema === 1 &&
      Array.isArray(stored.notes) &&
      Array.isArray(stored.proposals) &&
      Array.isArray(stored.branches) &&
      Array.isArray(stored.members) &&
      Array.isArray(stored.events)
    )
      s = { ...s, ...stored };
  } catch {}
  let observer;

  /* ------------------------------------------------------------------
     Local mode: served by `baker bakery`, which answers /api/memory with
     this project's real memory/ notes. Read-only. Without that endpoint
     (e.g. a static host) the page stays the sample workspace.
     ------------------------------------------------------------------ */
  let LOCAL = null,
    localSig = "";
  const KINDS = { decision: "Decision", lesson: "Lesson", knowledge: "Knowledge" };
  function wsName() {
    return LOCAL ? LOCAL.project : "Northstar";
  }
  function projectList() {
    return LOCAL ? [LOCAL.project] : ["Payments", "Storefront", "Platform"];
  }
  function prefsKey() {
    return "bakery-local-prefs:" + (LOCAL?.project || "");
  }
  function saveLocalPrefs() {
    try {
      localStorage.setItem(prefsKey(), JSON.stringify({ layout: s.layout, pins: s.pins, kind: s.kind }));
    } catch {}
  }
  function fmtDate(iso) {
    const m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(iso || "");
    if (!m) return iso || "unverified";
    const d = new Date(+m[1], +m[2] - 1, +m[3]),
      opts = { day: "numeric", month: "short" };
    if (d.getFullYear() !== new Date().getFullYear()) opts.year = "numeric";
    return d.toLocaleDateString("en-GB", opts);
  }
  function localNote(n, project) {
    const folder = n.path.split("/")[0].replace(/s$/, "");
    const kind = KINDS[(n.kind || "").toLowerCase()] || KINDS[folder] || "Knowledge";
    const firstPara = (n.body || "")
      .split(/\n\s*\n/)[0]
      .replace(/\[\[([^\]|#]+)[^\]]*\]\]/g, "$1")
      .replace(/\s+/g, " ")
      .trim();
    return {
      id: n.id,
      path: n.path,
      title: n.title,
      short: n.title.length > 46 ? n.title.slice(0, 44).trim() + "…" : n.title,
      kind,
      project,
      body: firstPara || "(no body yet)",
      fullBody: n.body || "",
      source: n.source || "no source recorded",
      asOf: n.as_of,
      lastVerified: n.last_verified,
      verified: fmtDate(n.last_verified),
      tags: n.tags || "",
      links: n.links || [],
      issues: n.issues || [],
      revision: "",
      author: "",
      initials: "",
    };
  }
  async function loadLocal(manual) {
    let snap;
    try {
      const res = await fetch("api/memory", { cache: "no-store" });
      if (!res.ok || !(res.headers.get("content-type") || "").includes("json")) throw new Error("no api");
      snap = await res.json();
      if (snap?.mode !== "local") throw new Error("not local");
    } catch {
      if (LOCAL) notify("Lost touch with baker. Is `baker bakery` still running?");
      return;
    }
    const sig = JSON.stringify([snap.notes, snap.findings, snap.branch, snap.head]);
    const first = !LOCAL;
    LOCAL = snap;
    if (first) enterLocal();
    if (sig === localSig && !first) {
      if (manual) notify("Up to date with memory/.");
      return;
    }
    localSig = sig;
    s.notes = snap.notes.map((n) => localNote(n, snap.project));
    s.pins = s.pins.filter((id) => s.notes.some((n) => n.id === id));
    if ($("#note-dialog").open) return; // don't yank content out from under a reader
    render();
    if (!first) notify(manual ? "Refreshed from memory/." : "memory/ changed — the table is updated.");
  }
  function enterLocal() {
    let prefs = {};
    try {
      prefs = JSON.parse(localStorage.getItem(prefsKey())) || {};
    } catch {}
    s = {
      ...seed(),
      view: ["memory", "connect"].includes(s.view) ? s.view : "memory",
      layout: prefs.layout || "cards",
      kind: prefs.kind || "All",
      pins: Array.isArray(prefs.pins) ? prefs.pins : [],
      project: "all",
      search: $("#global-search").value || "",
      proposals: [],
      branches: [],
      events: [],
      members: [],
    };
    document.body.classList.add("local-mode");
    document.title = LOCAL.project + " — Bakery (local)";
    $(".concept-bar").innerHTML =
      "<span>Baked locally.</span> Reading <code>memory/</code> from <strong>" +
      esc(LOCAL.project) +
      "</strong> · read-only · nothing leaves this machine " +
      '<button data-refresh-local="1">Refresh ↻</button>';
    $(".workspace-name").innerHTML =
      '<span class="workspace-symbol">' +
      esc(LOCAL.project.slice(0, 1).toLowerCase()) +
      '.</span><span>' +
      esc(LOCAL.project) +
      " <small>your project</small></span><span class=\"private-label\">local</span>";
    const nav = $(".sidebar nav");
    nav.querySelectorAll(".project-link").forEach((el) => el.remove());
    nav.insertAdjacentHTML(
      "beforeend",
      '<button class="nav-item project-link" data-project="' +
        esc(LOCAL.project) +
        '"><span class="project-dot payments"></span>' +
        esc(LOCAL.project) +
        "</button>",
    );
    $(".projects-label").textContent = "THIS PROJECT";
    $(".my-avatar").textContent = "you";
    $(".my-avatar").setAttribute("aria-label", "You, on this machine");
    $(".my-avatar").dataset.view = "connect";
    document.querySelectorAll('.nav-item[data-view="reviews"],.nav-item[data-view="branches"],.nav-item[data-view="activity"],.nav-item[data-view="team"]').forEach((el) => el.classList.add("demo-only"));
    if (location.hash && !/^#(memory|connect)/.test(location.hash)) history.replaceState(null, "", "#memory");
  }
  function localHealth() {
    if (!LOCAL.has_memory)
      return (
        '<div class="health-card warn">' +
        icon("book") +
        "<div><strong>No memory/ folder yet.</strong><p>Run <code>baker init</code> here, then let your agent leave notes as you work. They’ll show up on this table.</p></div></div>"
      );
    const issues = LOCAL.findings.filter((f) => f.level !== "ok");
    if (!issues.length)
      return (
        '<div class="health-card good">' +
        icon("check") +
        "<div><strong>" +
        esc((LOCAL.findings[0] || {}).msg || "Memory looks healthy.") +
        "</strong></div></div>"
      );
    return (
      '<details class="health-card warn"><summary>' +
      icon("clock") +
      "<strong>" +
      issues.length +
      " " +
      (issues.length === 1 ? "thing" : "things") +
      " to tidy in memory/</strong><span class=\"muted\">from <code>baker memory check</code></span></summary><ul>" +
      issues.map((f) => '<li class="' + esc(f.level) + '">' + esc(f.msg) + "</li>").join("") +
      "</ul></details>"
    );
  }
  function localPlaceholder(view) {
    const what = {
      reviewsView: ["Review queue", "Reviews"],
      branchesView: ["Branches", "Branches"],
      activityView: ["Activity", "Shared activity"],
      teamView: ["Team & access", "Team spaces"],
    }[view];
    return (
      '<div class="breadcrumb">' +
      esc(wsName()) +
      " <span>/</span> " +
      what[0] +
      "</div>" +
      head(
        "SOMETHING STILL IN THE OVEN.",
        what[1] + " <em>come with Bakery.</em>",
        "This local view only reads your project’s files. Shared review, branching, activity and team access need the hosted Bakery service, which isn’t built yet.",
      ) +
      '<div class="placeholder-paper"><img src="loaf.svg" alt="" width="70" height="74"><div><h2>Until then, Git does the sharing.</h2><p>Commit notes in <code>memory/</code> with your code and review them in pull requests. Teammates pull, and their agents read the same notes.</p>' +
      button("Back to memory →", 'data-view="memory"', "red") +
      "</div></div>"
    );
  }
  function openLocalNote(id) {
    const n = s.notes.find((x) => x.id === id);
    if (!n) return;
    s.focused = id;
    drawWeb();
    const related = neighbors(id);
    $("#note-content").innerHTML =
      '<header class="sheet-top"><span>' +
      icon(typeIcon(n.kind)) +
      esc(n.kind) +
      " · " +
      esc(n.project) +
      '</span><button class="icon-button" data-close="note-dialog" aria-label="Close crumb">×</button></header><div class="sheet-title"><div class="flex between">' +
      crumbStatus(n) +
      pinButton(n) +
      "</div><h2>" +
      esc(n.title) +
      '</h2><p class="mono">memory/' +
      esc(n.path) +
      "</p></div>" +
      (n.issues.length
        ? '<div class="warning-paper"><strong>Needs checking:</strong> ' +
          esc(n.issues.join("; ")) +
          ". Treat it as a lead and re-check the source before relying on it.</div>"
        : "") +
      '<section class="sheet-section"><span class="section-label">THE CRUMB</span><div class="note-body">' +
      esc(n.fullBody)
        .split(/\n\s*\n/)
        .map(
          (para) =>
            "<p>" +
            para
              .replace(/\n/g, "<br>")
              .replace(/\[\[([^\]|#]+)[^\]]*\]\]/g, (m, slug) => {
                const t = s.notes.find((x) => x.path.endsWith("/" + slug.trim() + ".md"));
                return t
                  ? '<button class="inline-link" data-open-note="' + esc(t.id) + '">' + esc(t.title) + "</button>"
                  : '<span class="broken-link" title="No note with this slug">' + slug + "</span>";
              }) +
            "</p>",
        )
        .join("") +
      '</div></section><section class="sheet-section"><h3>Follow the source</h3><div class="evidence-paper"><span>' +
      icon("book") +
      " " +
      esc(n.source) +
      "</span><p>As of " +
      esc(fmtDate(n.asOf)) +
      " · last verified " +
      esc(n.verified) +
      (n.tags ? " · tags " + esc(n.tags) : "") +
      '</p></div></section><section class="sheet-section"><h3>Connected crumbs <span class="muted">' +
      related.length +
      '</span></h3><div class="related-notes">' +
      (related.length
        ? related
            .map((o) => '<button data-open-note="' + esc(o.id) + '">' + icon("link") + "<span>" + esc(o.title) + "</span>↗</button>")
            .join("")
        : '<p class="muted">No links yet. Add <code>[[slug]]</code> in the note to connect it.</p>') +
      '</div></section><footer class="sheet-actions">' +
      button("Copy path", 'data-copy-path="memory/' + esc(n.path) + '"') +
      button("Copy Markdown", 'data-copy-note="' + esc(n.id) + '"') +
      "</footer>";
    if (!$("#note-dialog").open) $("#note-dialog").showModal();
    saveLocalPrefs();
  }
  let toastTimer;
  function notify(text) {
    clearTimeout(toastTimer);
    $("#toast").textContent = text;
    $("#toast").hidden = false;
    toastTimer = setTimeout(() => ($("#toast").hidden = true), 5200);
  }
  function save() {
    if (LOCAL) {
      saveLocalPrefs();
      return;
    }
    try {
      localStorage.setItem(KEY, JSON.stringify(s));
    } catch {
      notify("Browser storage is unavailable. Changes will last for this session.");
    }
  }
  function event(title, detail, kind = "update") {
    s.events.unshift({ title, detail, kind, when: "Just now", initials: "AM" });
  }
  function button(label, attrs = "", style = "") {
    return '<button class="button ' + style + '" ' + attrs + ">" + label + "</button>";
  }
  function tag(label, type = "") {
    return '<span class="tag ' + type + '">' + esc(label) + "</span>";
  }
  function head(eyebrow, title, description, actions = "") {
    return (
      '<div class="page-header"><div><span class="eyebrow">' +
      eyebrow +
      "</span><h1>" +
      title +
      "</h1><p>" +
      description +
      '</p></div><div class="page-actions">' +
      actions +
      "</div></div>"
    );
  }
  function typeIcon(kind) {
    return kind === "Decision" ? "shield" : kind === "Lesson" ? "spark" : "book";
  }
  function pending() {
    return s.proposals.filter((p) => p.status !== "merged");
  }
  function note(id) {
    return s.notes.find((n) => n.id === id);
  }
  function branch(name) {
    return s.branches.find((b) => b.name === name);
  }
  function proposal(id) {
    return s.proposals.find((p) => p.id === id);
  }
  function branchChanges(name) {
    return s.proposals.filter((p) => p.branch === name);
  }
  function neighbors(id) {
    const n = note(id);
    const ids = new Set(n?.links || []);
    s.notes.forEach((other) => {
      if (other.links.includes(id)) ids.add(other.id);
    });
    return [...ids].map(note).filter(Boolean);
  }
  function displayNotes() {
    let out = s.notes.map((n) => ({ ...n }));
    if (s.branch !== "main") {
      branchChanges(s.branch)
        .filter((p) => p.status !== "merged")
        .forEach((p) => {
          const i = out.findIndex((n) => n.id === p.noteId);
          const current = {
            ...(out[i] || {
              links: [],
              revision: 0,
              verified: "Draft",
              initials: "AM",
              commit: "Not merged",
            }),
            id: p.noteId,
            title: p.title,
            body: p.body,
            source: p.source,
            kind: p.kind,
            project: p.project,
            author: p.author,
            draft: p.status,
            proposalId: p.id,
          };
          if (i >= 0) out[i] = current;
          else out.push(current);
        });
    }
    return out.filter(
      (n) =>
        (s.project === "all" || n.project === s.project) &&
        (s.kind === "All" || (s.kind === "Pinned" && s.pins.includes(n.id)) || n.kind === s.kind) &&
        (!s.search ||
          (n.title + " " + n.body + " " + n.project + " " + n.source)
            .toLowerCase()
            .includes(s.search.toLowerCase())),
    );
  }
  function crumbStatus(n) {
    return n.draft
      ? tag(
          n.draft === "approved"
            ? "Reviewed draft"
            : n.draft === "changes"
              ? "Changes requested"
              : "Branch draft",
          "draft",
        )
      : n.issues?.length
        ? tag("Needs checking", "warning")
        : n.stale
          ? tag("Evidence changed", "warning")
          : tag(LOCAL ? "Fresh" : "Shared", "shared");
  }
  function pinButton(n) {
    const pinned = s.pins.includes(n.id);
    return (
      '<button class="pin-button ' +
      (pinned ? "pinned" : "") +
      '" data-pin="' +
      esc(n.id) +
      '" aria-pressed="' +
      pinned +
      '" aria-label="' +
      (pinned ? "Unpin" : "Pin") +
      " " +
      esc(n.title) +
      '">' +
      icon("pin") +
      "</button>"
    );
  }
  function crumbCard(n) {
    const update = s.proposals.some((p) => p.noteId === n.id && p.status !== "merged");
    return (
      '<article class="crumb-card" data-kind="' +
      n.kind +
      '" data-open-note="' +
      esc(n.id) +
      '"><div class="crumb-card-top"><span class="crumb-kind">' +
      icon(typeIcon(n.kind)) +
      n.kind +
      "</span>" +
      pinButton(n) +
      '</div><button class="crumb-title" data-open-note="' +
      esc(n.id) +
      '">' +
      esc(n.title) +
      "</button><p>" +
      esc(n.body) +
      '</p><div class="crumb-context"><span class="project-label">' +
      icon("folder") +
      esc(n.project) +
      "</span><span>" +
      neighbors(n.id).length +
      ' linked</span></div><div class="crumb-footer">' +
      crumbStatus(n) +
      (n.initials
        ? '<span class="avatar small" aria-label="' + esc(n.author) + '">' + esc(n.initials) + "</span>"
        : '<span class="footer-spacer"></span>') +
      '<span class="verified">' +
      (n.draft ? "Not on main" : "Verified " + esc(n.verified)) +
      "</span></div>" +
      (update && !n.draft ? '<span class="corner-note">↳ update on a branch</span>' : "") +
      "</article>"
    );
  }
  function crumbRow(n) {
    const update = s.proposals.some((p) => p.noteId === n.id && p.status !== "merged");
    return (
      '<article class="crumb-row" data-kind="' +
      n.kind +
      '" data-open-note="' +
      esc(n.id) +
      '"><span class="crumb-row-kind">' +
      icon(typeIcon(n.kind)) +
      '<span class="sr-only">' +
      n.kind +
      "</span></span>" +
      '<span class="crumb-row-main"><button class="crumb-row-title" data-open-note="' +
      esc(n.id) +
      '">' +
      esc(n.title) +
      "</button><span>" +
      esc(n.body) +
      "</span></span>" +
      '<span class="crumb-row-meta">' +
      esc(n.project) +
      "</span>" +
      '<span class="crumb-row-meta">' +
      neighbors(n.id).length +
      " linked</span>" +
      '<span class="crumb-row-status">' +
      crumbStatus(n) +
      (update && !n.draft ? '<span class="row-update">↳ update pending</span>' : "") +
      "</span>" +
      '<span class="crumb-row-meta">' +
      (n.draft ? "Not on main" : esc(n.verified)) +
      "</span>" +
      pinButton(n) +
      "</article>"
    );
  }
  // Deterministic force-directed layout so any number of crumbs gets its own spot.
  // Works in an aspect-corrected space (ASPECT wide × 1 tall) and returns percentages.
  const webLayoutCache = new Map();
  function webLayout(notes, aspect) {
    const key = aspect.toFixed(2) + "|" + notes.map((n) => n.id + ":" + n.links.join("+")).join(",");
    if (webLayoutCache.has(key)) return webLayoutCache.get(key);
    const N = notes.length,
      idx = new Map(notes.map((n, i) => [n.id, i])),
      projects = [...new Set(notes.map((n) => n.project))],
      order = notes.map((n, i) => i).sort((a, b) => projects.indexOf(notes[a].project) - projects.indexOf(notes[b].project) || a - b),
      P = new Array(N);
    order.forEach((i, k) => {
      const a = (2 * Math.PI * k) / Math.max(N, 1) - Math.PI / 2;
      P[i] = { x: aspect / 2 + 0.34 * aspect * Math.cos(a), y: 0.5 + 0.34 * Math.sin(a) };
    });
    const edges = [];
    notes.forEach((n) =>
      n.links.forEach((id) => {
        if (idx.has(id)) edges.push([idx.get(n.id), idx.get(id)]);
      }),
    );
    const W = Math.min(0.2 * aspect, 0.34),
      H = 0.2,
      minX = W / 2 + 0.03,
      maxX = aspect - W / 2 - 0.03,
      minY = 0.1 + H / 2,
      maxY = 0.92 - H / 2;
    for (let it = 0, t = 0.06; it < 400; it++, t *= 0.992) {
      const F = P.map(() => ({ x: 0, y: 0 }));
      for (let i = 0; i < N; i++)
        for (let j = i + 1; j < N; j++) {
          const dx = P[i].x - P[j].x,
            dy = P[i].y - P[j].y,
            d2 = Math.max(dx * dx + dy * dy, 1e-4),
            f = 0.012 / d2,
            d = Math.sqrt(d2);
          F[i].x += (dx / d) * f;
          F[i].y += (dy / d) * f;
          F[j].x -= (dx / d) * f;
          F[j].y -= (dy / d) * f;
        }
      edges.forEach(([i, j]) => {
        const dx = P[j].x - P[i].x,
          dy = P[j].y - P[i].y,
          d = Math.max(Math.hypot(dx, dy), 1e-3),
          f = (d - 0.42) * 0.6;
        F[i].x += (dx / d) * f;
        F[i].y += (dy / d) * f;
        F[j].x -= (dx / d) * f;
        F[j].y -= (dy / d) * f;
      });
      P.forEach((p, i) => {
        F[i].x += (aspect / 2 - p.x) * 0.05;
        F[i].y += (0.5 - p.y) * 0.05;
        const m = Math.hypot(F[i].x, F[i].y),
          step = Math.min(m, t);
        if (m > 0) {
          p.x = Math.min(maxX, Math.max(minX, p.x + (F[i].x / m) * step));
          p.y = Math.min(maxY, Math.max(minY, p.y + (F[i].y / m) * step));
        }
      });
    }
    // Final pass: push apart any overlapping note boxes.
    for (let it = 0; it < 120; it++) {
      let moved = false;
      for (let i = 0; i < N; i++)
        for (let j = i + 1; j < N; j++) {
          const ox = W - Math.abs(P[i].x - P[j].x),
            oy = H - Math.abs(P[i].y - P[j].y);
          if (ox > 0 && oy > 0) {
            moved = true;
            if (ox / aspect < oy) {
              const sx = (P[i].x < P[j].x ? -1 : 1) * (ox / 2 + 0.002);
              P[i].x = Math.min(maxX, Math.max(minX, P[i].x + sx));
              P[j].x = Math.min(maxX, Math.max(minX, P[j].x - sx));
            } else {
              const sy = (P[i].y < P[j].y ? -1 : 1) * (oy / 2 + 0.002);
              P[i].y = Math.min(maxY, Math.max(minY, P[i].y + sy));
              P[j].y = Math.min(maxY, Math.max(minY, P[j].y - sy));
            }
          }
        }
      if (!moved) break;
    }
    const out = P.map((p) => [(p.x / aspect) * 100, p.y * 100]);
    webLayoutCache.set(key, out);
    return out;
  }
  function memoryWeb(notes) {
    const width = $("#main")?.clientWidth || 900,
      height = Math.max(520, Math.min(860, 380 + notes.length * 34)),
      aspect = Math.max(0.6, Math.min(2.2, (width - 80) / height)),
      positions = webLayout(notes, aspect),
      legend = ["Decision", "Lesson", "Knowledge"]
        .map((k) => '<span class="legend-chip" data-kind="' + k + '">' + icon(typeIcon(k)) + k + "</span>")
        .join("");
    return (
      '<div class="memory-canvas" id="memory-canvas" style="height:' +
      height +
      'px"><div class="canvas-label"><span class="hand-star">✳</span> Hover a crumb to trace its connections.</div><div class="canvas-legend">' +
      legend +
      '</div><svg id="memory-edges" aria-hidden="true"></svg>' +
      notes
        .map((n, i) => {
          const [x, y] = positions[i];
          return (
            '<button class="web-node ' +
            (s.focused === n.id ? "selected" : "") +
            '" data-kind="' +
            n.kind +
            '" data-node="' +
            esc(n.id) +
            '" data-open-note="' +
            esc(n.id) +
            '" style="--x:' +
            x +
            "%;--y:" +
            y +
            '%"><span class="node-type">' +
            icon(typeIcon(n.kind)) +
            n.kind +
            "</span><strong>" +
            esc(n.short || n.title) +
            "</strong><span>" +
            esc(n.project) +
            " · " +
            neighbors(n.id).length +
            " linked" +
            (n.stale ? " · needs checking" : "") +
            "</span></button>"
          );
        })
        .join("") +
      '<div class="canvas-bottom"><span>Each line is an explicit note link.</span><span>' +
      notes.length +
      " crumbs on the table</span></div></div>"
    );
  }
  let tracedId = "";
  function traceWeb(id) {
    tracedId = id;
    const canvas = $("#memory-canvas");
    if (!canvas) return;
    const related = new Set(id ? neighbors(id).map((n) => n.id) : []);
    canvas.classList.toggle("tracing", !!id);
    canvas.querySelectorAll("[data-node]").forEach((el) => {
      el.classList.toggle("traced", el.dataset.node === id);
      el.classList.toggle("related", related.has(el.dataset.node));
    });
    canvas.querySelectorAll("#memory-edges path").forEach((p) => {
      p.classList.toggle("traced", !!id && (p.dataset.a === id || p.dataset.b === id));
    });
  }
  function memoryView() {
    const notes = displayNotes();
    const isBranch = s.branch !== "main";
    return (
      '<div class="breadcrumb">' + esc(wsName()) + ' <span>/</span> ' +
      (s.project === "all" ? "The shared table" : esc(s.project)) +
      "</div>" +
      head(
        "LITTLE DISCOVERIES. SHARED UNDERSTANDING.",
        isBranch
          ? "A little room to <em>try things.</em>"
          : 'Good ideas,<br class="desktop-break"> <em>kept together.</em>',
        isBranch
          ? "You’re looking at a branch. Main keeps its current version until changes are reviewed and merged."
          : "Decisions, lessons, and the details you’d rather not explain twice.",
        LOCAL ? button("+ Add a crumb", "data-local-howto", "red") : button("+ Leave a crumb", "data-compose", "red"),
      ) +
      (LOCAL ? localHealth() : "") +
      '<div class="memory-toolbar"><div class="memory-filters" aria-label="Memory kind">' +
      ["All", "Decision", "Lesson", "Knowledge", "Pinned"]
        .map(
          (k) =>
            '<button class="' +
            (s.kind === k ? "selected" : "") +
            '" data-kind-filter="' +
            k +
            '" aria-pressed="' +
            (s.kind === k) +
            '">' +
            (k === "Pinned" ? icon("pin") : "") +
            (k === "All" ? "All crumbs" : k === "Decision" ? "Decisions" : k === "Lesson" ? "Lessons" : k) +
            "</button>",
        )
        .join("") +
      '</div><div class="view-switch" aria-label="Memory layout">' +
      [
        ["cards", "grid", "Cards"],
        ["list", "notes", "List"],
        ["web", "network", "Web"],
      ]
        .map(
          ([key, i, label]) =>
            '<button class="' +
            (s.layout === key ? "selected" : "") +
            '" data-layout="' +
            key +
            '" aria-pressed="' +
            (s.layout === key) +
            '">' +
            icon(i) +
            label +
            "</button>",
        )
        .join("") +
      '</div></div><div class="context-toolbar">' +
      (LOCAL
        ? '<span class="branch-select local-branch">' +
          icon("branch") +
          esc(LOCAL.branch || "working tree") +
          (LOCAL.head ? " · " + esc(LOCAL.head) : "") +
          ' · <span class="muted">files on disk, including uncommitted edits</span></span><label hidden>'
        : '<label class="branch-select">') +
      icon("branch") +
      '<span class="sr-only">Memory branch</span><select id="memory-branch"><option value="main" ' +
      (!isBranch ? "selected" : "") +
      ">main · shared memory</option>" +
      s.branches
        .filter((b) => b.status === "open")
        .map(
          (b) =>
            '<option value="' +
            esc(b.name) +
            '" ' +
            (s.branch === b.name ? "selected" : "") +
            ">" +
            esc(b.name) +
            "</option>",
        )
        .join("") +
      '</select></label><label class="project-select"><span class="sr-only">Filter by project</span><select id="memory-project">' +
      ["all", ...projectList()]
        .map(
          (p) =>
            '<option value="' +
            p +
            '" ' +
            (s.project === p ? "selected" : "") +
            ">" +
            (p === "all" ? "All projects" : p) +
            "</option>",
        )
        .join("") +
      '</select></label><span class="context-count">' +
      notes.length +
      " " +
      (notes.length === 1 ? "crumb" : "crumbs") +
      "</span></div>" +
      (s.search
        ? '<div class="search-notice">Results for “' +
          esc(s.search) +
          "” " +
          button("Clear search", "data-clear-search", "quiet") +
          "</div>"
        : "") +
      (isBranch
        ? '<div class="branch-notice">' +
          icon("branch") +
          "<span>" +
          esc(s.branch) +
          " · changes stay separate from main</span>" +
          button("Open branch →", 'data-open-branch="' + esc(s.branch) + '"', "quiet") +
          "</div>"
        : "") +
      (notes.length
        ? s.layout === "web"
          ? memoryWeb(notes)
          : s.layout === "list"
            ? '<div class="crumb-list">' + notes.map(crumbRow).join("") + "</div>"
            : '<div class="crumb-grid">' + notes.map(crumbCard).join("") + "</div>"
        : '<div class="empty-state"><img src="loaf.svg" alt="" width="90" height="95"><h2>No crumbs on this part of the table.</h2><p>Try another filter, or leave a useful first note.</p>' +
          button("Show all crumbs", "data-clear-filters") +
          "</div>") +
      '<div class="memory-bottom-note"><span>✳</span><p>Small notes. Clear sources. A little less starting over.</p><button data-view="connect">Your memory is still yours ↗</button></div>'
    );
  }
  function branchState(name) {
    const changes = branchChanges(name),
      approved = changes.filter((p) => p.status === "approved" || p.status === "merged").length,
      conflict = changes.some(
        (p) => note(p.noteId) && note(p.noteId).revision !== p.baseRevision && p.status !== "merged",
      ),
      merged = branch(name)?.status === "merged";
    return {
      changes,
      approved,
      total: changes.length,
      conflict,
      merged,
      ready: !merged && !conflict && changes.length > 0 && approved === changes.length,
    };
  }
  function nextStep(p) {
    if (p.status === "merged") return "On main";
    const b = branchState(p.branch);
    if (b.conflict) return "Main moved · update branch";
    if (p.status === "changes") return "Waiting on " + p.author.split(" ")[0];
    if (p.status === "pending") return "Waiting on a reviewer";
    if (b.ready) return "Ready to merge";
    return b.approved + " of " + b.total + " approved on branch";
  }
  function statusTag(p) {
    return tag(
      p.status === "pending"
        ? "Needs review"
        : p.status === "changes"
          ? "Changes requested"
          : p.status === "approved"
            ? "Approved"
            : "Merged",
      p.status === "merged" || p.status === "approved" ? "shared" : p.status === "changes" ? "draft" : "warning",
    );
  }
  function proposalRow(p) {
    return (
      '<button class="proposal-row" data-review="' +
      esc(p.id) +
      '"><span class="review-icon">' +
      icon("review") +
      '</span><span class="proposal-row-content"><strong>' +
      esc(p.title) +
      "</strong><span>" +
      esc(p.project) +
      " · " +
      esc(p.author) +
      " via " +
      esc(p.agent || "Bakery") +
      " · <code>" +
      esc(p.branch) +
      '</code></span></span><span class="proposal-row-status">' +
      statusTag(p) +
      '<span class="next-step">' +
      esc(nextStep(p)) +
      "</span></span>" +
      '<span class="row-arrow">↗</span></button>'
    );
  }
  function reviewSummary() {
    const open = s.proposals.filter((p) => p.status !== "merged"),
      count = (st) => open.filter((p) => p.status === st).length,
      ready = s.branches.filter((b) => b.status === "open" && branchState(b.name).ready),
      stale = s.branches.filter((b) => b.status === "open" && branchState(b.name).conflict);
    const stat = (n, label, cls = "") =>
      '<div class="review-stat ' + cls + '"><strong>' + n + "</strong><span>" + label + "</span></div>";
    return (
      '<div class="review-summary">' +
      stat(count("pending"), "need review", count("pending") ? "attention" : "") +
      stat(count("changes"), "waiting on authors") +
      stat(count("approved"), "approved") +
      stat(ready.length, ready.length === 1 ? "branch ready to merge" : "branches ready to merge", ready.length ? "good" : "") +
      "</div>" +
      (ready.length || stale.length
        ? '<div class="review-callouts">' +
          ready
            .map(
              (b) =>
                '<div class="callout good">' +
                icon("check") +
                "<span><code>" +
                esc(b.name) +
                "</code> is reviewed and ready.</span>" +
                button("Merge →", 'data-merge="' + esc(b.name) + '"', "red") +
                "</div>",
            )
            .join("") +
          stale
            .map(
              (b) =>
                '<div class="callout warn">' +
                icon("branch") +
                "<span><code>" +
                esc(b.name) +
                "</code> started before main changed.</span>" +
                button("Update from main", 'data-rebase="' + esc(b.name) + '"') +
                "</div>",
            )
            .join("") +
          "</div>"
        : "")
    );
  }
  function reviewsView() {
    if (LOCAL) return localPlaceholder('reviewsView');
    if (s.selectedReview && proposal(s.selectedReview)) return reviewDetail(proposal(s.selectedReview));
    const rows = s.proposals.filter((p) =>
      s.reviewFilter === "merged" ? p.status === "merged" : p.status !== "merged",
    );
    return (
      '<div class="breadcrumb">' + esc(wsName()) + ' <span>/</span> Review queue</div>' +
      head(
        "A SECOND PAIR OF EYES.",
        "Good knowledge<br> <em>deserves a good look.</em>",
        "Check the claim, its context, and its evidence before it joins the shared memory.",
      ) +
      (s.reviewFilter === "open" ? reviewSummary() : "") +
      '<div class="simple-tabs">' +
      ["open", "merged"]
        .map(
          (f) =>
            '<button class="' +
            (s.reviewFilter === f ? "selected" : "") +
            '" data-review-filter="' +
            f +
            '">' +
            (f === "open" ? "On the counter" : "Put away") +
            " <span>" +
            s.proposals.filter((p) => (f === "merged" ? p.status === "merged" : p.status !== "merged"))
              .length +
            "</span></button>",
        )
        .join("") +
      '</div><div class="proposal-list">' +
      (rows.length
        ? rows.map(proposalRow).join("")
        : '<div class="empty-state"><h2>All caught up.</h2><p>A clean counter. A lovely thing.</p></div>') +
      "</div>"
    );
  }
  function reviewNextStep(p) {
    const b = branchState(p.branch),
      bar =
        '<div class="approval-progress" aria-label="' +
        b.approved +
        " of " +
        b.total +
        ' changes approved"><span style="width:' +
        (b.total ? Math.round((b.approved / b.total) * 100) : 0) +
        '%"></span></div>',
      line = (title, text, action = "", cls = "") =>
        '<div class="next-step-card ' +
        cls +
        '"><span class="section-label">NEXT STEP</span><strong>' +
        title +
        "</strong><p>" +
        text +
        "</p>" +
        action +
        "</div>";
    if (b.merged || p.status === "merged")
      return line("Merged into main.", "This change is part of the shared memory.", "", "good");
    if (b.conflict)
      return line(
        "Update the branch first.",
        "Main changed after this draft began. Updating keeps the draft, refreshes its base, and clears approvals.",
        button("Update from main", 'data-rebase="' + esc(p.branch) + '"', "full"),
        "warn",
      );
    const progress =
      bar + '<p class="progress-label">' + b.approved + " of " + b.total + " changes on this branch approved</p>";
    if (b.ready)
      return line(
        "Ready to merge.",
        "Every change on <code>" + esc(p.branch) + "</code> is approved.",
        progress + button("Merge branch →", 'data-merge="' + esc(p.branch) + '"', "red full"),
        "good",
      );
    if (p.status === "approved")
      return line("Approved. Waiting on the rest of the branch.", "The branch merges once all its changes are approved.", progress);
    if (p.status === "changes")
      return line("Waiting on " + esc(p.author.split(" ")[0]) + ".", "Changes were requested. Main keeps its existing version.", progress);
    return line("Your review.", "Approving doesn’t merge. Main stays as it is until the branch is merged.", progress);
  }
  function reviewDetail(p) {
    const current = note(p.noteId);
    const isMerged = p.status === "merged";
    const changed = current && current.revision !== p.baseRevision && !isMerged;
    return (
      '<div class="breadcrumb"><button data-back-reviews>← Review queue</button><span>/</span>' +
      esc(p.project) +
      '</div><div class="review-title"><span class="eyebrow">' +
      icon("branch") +
      esc(p.branch) +
      "</span><h1>" +
      esc(p.title) +
      '</h1><div class="review-byline"><span class="avatar">' +
      esc(
        p.author
          .split(" ")
          .map((x) => x[0])
          .join(""),
      ) +
      "</span><span>Proposed by " +
      esc(p.author) +
      ' <span class="muted">via ' +
      esc(p.agent || "Bakery") +
      "</span></span>" +
      tag(
        p.status === "pending"
          ? "Needs review"
          : p.status === "changes"
            ? "Changes requested"
            : p.status === "approved"
              ? "Approved"
              : "Merged",
        p.status === "approved" || isMerged ? "shared" : "warning",
      ) +
      '</div></div><div class="review-body-layout"><div><section class="review-section"><span class="section-label">WHY THIS CHANGE</span><p>' +
      esc(p.reason) +
      '</p></section><section class="review-section"><div class="section-heading"><h2>The change</h2><span class="mono muted">' +
      (p.baseRevision ? "revision " + p.baseRevision + " → " + (p.baseRevision + 1) : "a new crumb") +
      '</span></div><div class="diff-paper"><div class="diff-label">' +
      (p.baseRevision ? "CURRENTLY SHARED" : "NOT YET SHARED") +
      '</div><div class="diff-old">− ' +
      esc(p.baseBody || "No previous note.") +
      '</div><div class="diff-label">ON THIS BRANCH</div><div class="diff-new">+ ' +
      esc(p.body) +
      '</div></div></section><section class="review-section"><h2>Evidence to look at</h2><div class="evidence-paper"><span>' +
      icon("book") +
      " " +
      esc(p.source) +
      "</span><p>Sample reference supplied with this proposal. A reviewer should inspect the real source before relying on the claim.</p>" +
      (current
        ? "<details><summary>See the sample evidence</summary><pre>" +
          esc(current.evidence) +
          "</pre></details>"
        : "") +
      '</div></section><section class="review-section"><h2>A little discussion</h2><div class="discussion">' +
      p.comments
        .map(
          (c) =>
            '<div class="comment"><span class="avatar small">' +
            esc(c.author.slice(0, 2).toUpperCase()) +
            "</span><div><strong>" +
            esc(c.author) +
            "</strong><p>" +
            esc(c.text) +
            "</p><small>" +
            esc(c.when) +
            "</small></div></div>",
        )
        .join("") +
      "</div>" +
      (!isMerged
        ? '<form class="review-comment-form" data-review-form="' +
          esc(p.id) +
          '"><label class="sr-only" for="review-comment">Review comment</label><textarea id="review-comment" name="comment" rows="3" placeholder="Ask a question, or leave a useful reason…"></textarea><div class="review-actions"><span>Reviewing as Alex · sample owner</span><div>' +
          button("Comment", 'type="submit" name="action" value="comment"') +
          button("Request changes", 'type="submit" name="action" value="changes"') +
          button(
            p.status === "approved" ? "Approved ✓" : "Approve change",
            'type="submit" name="action" value="approve" ' +
              (p.status === "approved" || changed ? "disabled" : ""),
            "red",
          ) +
          "</div></div></form>"
        : "") +
      '</section></div><aside class="review-aside"><div class="review-aside-card"><span class="section-label">A PLACE IN THE NOTEBOOK</span><dl><dt>Project</dt><dd>' +
      esc(p.project) +
      "</dd><dt>Kind</dt><dd>" +
      esc(p.kind) +
      "</dd><dt>Target</dt><dd>main · shared memory</dd><dt>Visibility</dt><dd>Private · Northstar</dd></dl>" +
      reviewNextStep(p) +
      button("Open branch →", 'data-open-branch="' + esc(p.branch) + '"', "full quiet-full") +
      "</div>" +
      "</aside></div>"
    );
  }
  function branchesView() {
    if (LOCAL) return localPlaceholder('branchesView');
    if (s.selectedBranch && branch(s.selectedBranch)) return branchDetail(branch(s.selectedBranch));
    return (
      '<div class="breadcrumb">' + esc(wsName()) + ' <span>/</span> Branches</div>' +
      head(
        "SAME NOTEBOOK. ROOM TO EXPERIMENT.",
        "Try a different<br> <em>recipe.</em>",
        "Keep a new idea separate while the team works out the details.",
        button("+ Make a branch", "data-new-branch", "red"),
      ) +
      '<div class="main-branch-card"><span class="branch-symbol">' +
      icon("branch") +
      '</span><div><h2>main <span class="tag shared">Shared memory</span></h2><p>The version your team’s agents should start from.</p></div>' +
      button("Browse main →", "data-browse-main") +
      '</div><div class="section-heading branch-list-title"><h2>On the workbench</h2><span class="muted small">' +
      s.branches.filter((b) => b.status === "open").length +
      ' open branches</span></div><div class="branch-grid">' +
      s.branches
        .map((b) => {
          const changes = branchChanges(b.name),
            ready = changes.length && changes.every((p) => p.status === "approved" || p.status === "merged");
          return (
            '<button class="branch-card" data-open-branch="' +
            esc(b.name) +
            '"><div class="branch-card-top">' +
            icon("branch") +
            tag(
              b.status === "merged" ? "Merged" : ready ? "Ready to merge" : "In progress",
              b.status === "merged" || ready ? "shared" : "",
            ) +
            "</div><h3>" +
            esc(b.name) +
            "</h3><p>" +
            esc(b.description || "A fresh place for a useful change.") +
            "</p><footer><span>" +
            esc(b.project) +
            "</span><span>" +
            changes.length +
            " " +
            (changes.length === 1 ? "change" : "changes") +
            " ↗</span></footer></button>"
          );
        })
        .join("") +
      "</div>"
    );
  }
  function branchDetail(b) {
    const changes = branchChanges(b.name);
    const ready = changes.length > 0 && changes.every((p) => p.status === "approved");
    const conflict = changes.some(
      (p) => note(p.noteId) && note(p.noteId).revision !== p.baseRevision && p.status !== "merged",
    );
    return (
      '<div class="breadcrumb"><button data-back-branches>← Branches</button><span>/</span>' +
      esc(b.project) +
      "</div>" +
      head(
        "A WORKING COPY, WITH ROOM TO THINK.",
        esc(b.name),
        "Changes stay separate from the shared notebook until the branch is merged.",
        button("Browse branch", 'data-browse-branch="' + esc(b.name) + '"'),
      ) +
      '<div class="branch-summary"><span>' +
      icon("branch") +
      " " +
      esc(b.name) +
      '</span><span class="branch-line" aria-hidden="true">→</span><span>' +
      icon("notes") +
      " main</span>" +
      tag(
        b.status === "merged" ? "Merged" : ready ? "Reviewed and ready" : "Review in progress",
        b.status === "merged" || ready ? "shared" : "warning",
      ) +
      (b.status !== "merged" && changes.length
        ? '<span class="branch-progress">' +
          branchState(b.name).approved +
          " of " +
          changes.length +
          " approved</span>"
        : "") +
      "</div>" +
      (conflict
        ? '<div class="warning-paper"><p>Shared memory changed after this branch began. Updating from main keeps your draft, refreshes its base, and clears previous approvals.</p>' +
          button("Update from main", 'data-rebase="' + esc(b.name) + '"') +
          "</div>"
        : "") +
      '<div class="proposal-list">' +
      (changes.length
        ? changes.map(proposalRow).join("")
        : '<div class="empty-state"><h2>A clean sheet of paper.</h2><p>Leave the first crumb on this branch.</p>' +
          button("+ Leave a crumb", 'data-compose data-compose-branch="' + esc(b.name) + '"', "red") +
          "</div>") +
      '</div><div class="merge-paper"><div><h2>' +
      (b.status === "merged"
        ? "This batch is on the shared table."
        : ready
          ? "Ready to join the notebook."
          : "A little review before the merge.") +
      "</h2><p>" +
      (b.status === "merged"
        ? "These changes are now part of main."
        : ready
          ? "All changes have been approved. Merge to update the shared memory."
          : "Every proposed change needs approval before this branch can merge.") +
      "</p></div>" +
      button(
        b.status === "merged" ? "Merged ✓" : "Merge branch →",
        'data-merge="' + esc(b.name) + '" ' + (!ready || conflict || b.status === "merged" ? "disabled" : ""),
        "red",
      ) +
      "</div>"
    );
  }
  function activityView() {
    if (LOCAL) return localPlaceholder('activityView');
    return (
      '<div class="breadcrumb">' + esc(wsName()) + ' <span>/</span> Activity</div>' +
      head(
        "A TRAIL OF GOOD THINKING.",
        "What’s been<br> <em>cooking?</em>",
        "The people, changes, and decisions behind the shared memory.",
      ) +
      '<div class="activity-timeline">' +
      s.events
        .map(
          (e) =>
            '<article class="activity-item"><span class="avatar">' +
            esc(e.initials) +
            "</span><div><h3>" +
            esc(e.title) +
            "</h3><p>" +
            esc(e.detail) +
            "</p><time>" +
            esc(e.when) +
            '</time></div><span class="activity-kind">' +
            icon(e.kind === "merge" ? "check" : e.kind === "proposal" ? "review" : "clock") +
            "</span></article>",
        )
        .join("") +
      "</div>"
    );
  }
  function teamView() {
    if (LOCAL) return localPlaceholder('teamView');
    return (
      '<div class="breadcrumb">' + esc(wsName()) + ' <span>/</span> Team & access</div>' +
      head(
        "THERE’S ROOM AT THE TABLE.",
        "Good company.<br> <em>Shared understanding.</em>",
        "People own the decisions. Agents help keep the useful context.",
        button("+ Invite a teammate", "data-invite", "red"),
      ) +
      '<div class="team-paper"><div class="team-table-head"><span>Teammate</span><span>Role</span></div>' +
      s.members
        .map(
          (m) =>
            '<div class="member-row"><span class="avatar">' +
            esc(m.initials) +
            "</span><div><strong>" +
            esc(m.name) +
            "</strong><span>" +
            esc(m.email) +
            "</span></div>" +
            tag(m.role, m.role === "Owner" ? "shared" : "") +
            "</div>",
        )
        .join("") +
      '</div><div class="role-notes"><article><h3>Owner</h3><p>Manage workspace membership and policy.</p></article><article><h3>Reviewer</h3><p>Review proposed knowledge and approve changes.</p></article><article><h3>Contributor</h3><p>Use shared memory and propose changes on branches.</p></article></div><div class="small-print">Roles here illustrate the proposed experience. This concept does not implement authentication or enforce team permissions.</div>'
    );
  }
  function connectView() {
    return (
      '<div class="breadcrumb">' + esc(wsName()) + ' <span>/</span> Baker & your files</div>' +
      head(
        "BAKED LOCALLY. YOURS TO KEEP.",
        "Your crumbs.<br> <em>Your cupboard.</em>",
        "Baker works on its own. Git shares the files. Bakery makes room for a better team experience.",
      ) +
      '<div class="ownership-layout"><div class="ownership-paper"><span class="section-label">THE PORTABLE PART</span><pre>memory/\n├── MEMORY.md\n├── decisions/\n├── lessons/\n└── knowledge/</pre><p>Small Markdown files, with sources and dates. Readable in an editor. Shareable through Git.</p>' +
      button(icon("download") + " Export shared memory", "data-export", "red") +
      '<small>Downloads the sample main-branch notes as Markdown.</small></div><div class="ownership-notes"><article><span>01</span><div><h2>Baker remembers locally.</h2><p>Retain notes, checkpoint work, and switch agents without a Bakery account.</p></div></article><article><span>02</span><div><h2>Git already shares.</h2><p>Commit memory with the code. Teammates can pull, branch, and review it with their existing tools.</p></div></article><article><span>03</span><div><h2>Bakery adds the shared table.</h2><p>Visual connections, memory-specific review, and room to understand how a decision changed.</p></div></article><a class="text-link" href="https://baker-fresh-memory.chatgpt-for-3601.chatgpt.site/docs/teams/" target="_blank" rel="noopener noreferrer">Read the Baker team workflow ↗</a></div></div><aside class="local-boundary"><img src="loaf.svg" width="76" height="80" alt=""><div><h3>The personal bits stay personal.</h3><p>Handoff notes, raw transcripts, and local checkpoints are not automatically team memory. Share useful knowledge deliberately.</p></div></aside>'
    );
  }
  function fillIcons() {
    document.querySelectorAll("[data-icon]").forEach((el) => (el.innerHTML = icon(el.dataset.icon)));
  }
  function render() {
    if (observer) {
      observer.disconnect();
      observer = null;
    }
    $("#main").innerHTML =
      s.view === "reviews"
        ? reviewsView()
        : s.view === "branches"
          ? branchesView()
          : s.view === "activity"
            ? activityView()
            : s.view === "team"
              ? teamView()
              : s.view === "connect"
                ? connectView()
                : memoryView();
    $("#review-count").textContent = pending().length;
    document
      .querySelectorAll(".nav-item[data-view]")
      .forEach((el) => el.classList.toggle("active", el.dataset.view === s.view));
    document
      .querySelectorAll("[data-project]")
      .forEach((el) =>
        el.classList.toggle("project-active", s.view === "memory" && s.project === el.dataset.project),
      );
    fillIcons();
    mountWeb();
  }
  function go(view, id = "") {
    s.view = view;
    s.selectedReview = view === "reviews" ? id : "";
    s.selectedBranch = view === "branches" ? id : "";
    history.pushState(null, "", "#" + view + (id ? "/" + encodeURIComponent(id) : ""));
    render();
    save();
    window.scrollTo({ top: 0, behavior: "instant" });
  }
  function route() {
    const [view, id] = location.hash.slice(1).split("/");
    if (["memory", "reviews", "branches", "activity", "team", "connect"].includes(view)) {
      s.view = view;
      try {
        s.selectedReview = view === "reviews" ? decodeURIComponent(id || "") : "";
        s.selectedBranch = view === "branches" ? decodeURIComponent(id || "") : "";
      } catch {
        s.selectedReview = "";
        s.selectedBranch = "";
      }
    }
    render();
  }
  function openNote(id) {
    if (LOCAL) return openLocalNote(id);
    const n = displayNotes().find((x) => x.id === id) || note(id);
    if (!n) return;
    s.focused = id;
    drawWeb();
    const related = neighbors(id);
    const update = s.proposals.find((p) => p.noteId === id && p.status !== "merged");
    $("#note-content").innerHTML =
      '<header class="sheet-top"><span>' +
      icon(typeIcon(n.kind)) +
      esc(n.kind) +
      " · " +
      esc(n.project) +
      '</span><button class="icon-button" data-close="note-dialog" aria-label="Close crumb">×</button></header><div class="sheet-title"><div class="flex between">' +
      crumbStatus(n) +
      pinButton(n) +
      "</div><h2>" +
      esc(n.title) +
      "</h2><p>" +
      (n.draft ? "Based on revision " : "Revision ") +
      n.revision +
      " · " +
      esc(n.author) +
      "</p></div>" +
      (n.stale
        ? '<div class="warning-paper">The referenced evidence changed after verification. Recheck the current source before relying on this note.</div>'
        : "") +
      (n.draft
        ? '<div class="branch-notice">' + icon("branch") + " This is a draft on " + esc(s.branch) + ".</div>"
        : "") +
      '<section class="sheet-section"><span class="section-label">THE CRUMB</span><p class="note-body">' +
      esc(n.body).replace(/\n/g, "<br>") +
      '</p></section><section class="sheet-section"><h3>Why we keep it</h3><p>' +
      esc(n.why || "A useful piece of context for the next developer and agent.") +
      '</p></section><section class="sheet-section"><h3>Follow the source</h3><div class="evidence-paper"><span>' +
      icon("book") +
      " " +
      esc(n.source) +
      "</span><p>" +
      (n.draft
        ? "Proposed source reference · not yet reviewed"
        : "Recorded at " + esc(n.commit) + " · verified " + esc(n.verified)) +
      "</p><details><summary>Peek at the sample evidence</summary><pre>" +
      esc(n.evidence || "Source reference supplied in the branch proposal.") +
      '</pre></details></div></section><section class="sheet-section"><h3>Connected crumbs <span class="muted">' +
      related.length +
      '</span></h3><div class="related-notes">' +
      (related.length
        ? related
            .map(
              (o) =>
                '<button data-open-note="' +
                esc(o.id) +
                '">' +
                icon("link") +
                "<span>" +
                esc(o.title) +
                "</span>↗</button>",
            )
            .join("")
        : '<p class="muted">No links yet. A little room to connect the dots.</p>') +
      '</div></section><details class="sheet-history"><summary>How this crumb got here</summary><p>Revision ' +
      n.revision +
      " " +
      (n.draft ? "is on a branch." : "is shared on main.") +
      " " +
      esc(n.author) +
      " retained this context from " +
      esc(n.source) +
      ".</p></details>" +
      (update
        ? '<div class="update-note"><span>↳ There’s an update on a branch.</span>' +
          button("See the change", 'data-review="' + esc(update.id) + '"', "quiet") +
          "</div>"
        : "") +
      '<footer class="sheet-actions">' +
      button("Copy Markdown", 'data-copy-note="' + esc(n.id) + '"') +
      button(icon("edit") + " Propose an edit", 'data-edit-note="' + esc(n.id) + '"', "red") +
      "</footer>";
    if (!$("#note-dialog").open) $("#note-dialog").showModal();
    save();
  }
  function composer(id = "", branchName = "") {
    const form = $("#crumb-form");
    form.reset();
    $("#crumb-error").textContent = "";
    const n = id ? displayNotes().find((x) => x.id === id) || note(id) : null;
    form.elements.editId.value = id;
    form.elements.project.value = n?.project || (s.project === "all" ? "Payments" : s.project);
    form.elements.kind.value = n?.kind || "Decision";
    form.elements.title.value = n?.title || "";
    form.elements.body.value = n?.body || "";
    form.elements.source.value = n?.source || "";
    form.elements.branch.value =
      branchName || (s.branch !== "main" ? s.branch : id ? "notes/update-" + id : "");
    $("#composer-title").textContent = id ? "A better version." : "Leave a crumb.";
    $("#note-dialog").close();
    $("#composer-dialog").showModal();
  }
  function markdown(n) {
    if (n.path)
      return (
        "---\ntype: " + n.kind.toLowerCase() + "\n" +
        (n.asOf ? "as_of: " + n.asOf + "\n" : "") +
        (n.lastVerified ? "last_verified: " + n.lastVerified + "\n" : "") +
        "source: " + n.source + "\n---\n# " + n.title + "\n\n" + n.fullBody + "\n"
      );
    const match = String(n.verified).match(/^(\d+) (Nov|Oct)$/);
    const verified = match
      ? "2026-" + (match[2] === "Nov" ? "11" : "10") + "-" + match[1].padStart(2, "0")
      : null;
    return (
      "---\ntype: " +
      n.kind.toLowerCase() +
      "\nas_of: 2026-11-05\n" +
      (verified ? "last_verified: " + verified + "\n" : "") +
      "source: " +
      JSON.stringify(n.source) +
      "\n---\n# " +
      n.title +
      "\n\n" +
      n.body +
      "\n" +
      (n.links?.length ? "\n" + n.links.map((id) => "[[" + id + "]]").join(" · ") + "\n" : "")
    );
  }
  async function copyNote(id) {
    const n = displayNotes().find((x) => x.id === id) || note(id);
    try {
      await navigator.clipboard.writeText(markdown(n));
      notify("Crumb copied as Markdown. Yours to take with you.");
    } catch {
      notify("Clipboard unavailable. Export the shared notes from Baker & your files.");
    }
  }
  function drawWeb() {
    const canvas = $("#memory-canvas");
    if (!canvas) return;
    const svg = $("#memory-edges"),
      rect = canvas.getBoundingClientRect();
    svg.setAttribute("viewBox", "0 0 " + canvas.clientWidth + " " + canvas.clientHeight);
    const seen = new Set(),
      edges = [];
    canvas.querySelectorAll("[data-node]").forEach((el) => {
      const n = note(el.dataset.node);
      if (!n) return;
      n.links.forEach((id) => {
        const to = canvas.querySelector('[data-node="' + id + '"]');
        if (!to) return;
        const key = [n.id, id].sort().join(":");
        if (seen.has(key)) return;
        seen.add(key);
        const a = el.getBoundingClientRect(),
          b = to.getBoundingClientRect(),
          x1 = a.left + a.width / 2 - rect.left,
          y1 = a.top + a.height / 2 - rect.top,
          x2 = b.left + b.width / 2 - rect.left,
          y2 = b.top + b.height / 2 - rect.top;
        edges.push(
          '<path data-a="' +
            esc(n.id) +
            '" data-b="' +
            esc(id) +
            '" class="' +
            (s.focused === n.id || s.focused === id ? "focused" : "") +
            '" d="M ' +
            x1 +
            " " +
            y1 +
            " C " +
            x1 +
            " " +
            (y1 + y2) / 2 +
            " " +
            x2 +
            " " +
            (y1 + y2) / 2 +
            " " +
            x2 +
            " " +
            y2 +
            '"/>',
        );
      });
    });
    svg.innerHTML = edges.join("");
    if (tracedId) traceWeb(tracedId);
  }
  function mountWeb() {
    const canvas = $("#memory-canvas");
    if (canvas) {
      drawWeb();
      observer = new ResizeObserver(drawWeb);
      observer.observe(canvas);
      canvas.querySelectorAll("[data-node]").forEach((n) => observer.observe(n));
    }
  }
  function merge(name) {
    const b = branch(name),
      changes = branchChanges(name);
    if (!b || b.status === "merged" || !changes.length || changes.some((p) => p.status !== "approved")) {
      notify("Every change needs review before this branch can merge.");
      return;
    }
    if (changes.some((p) => note(p.noteId) && note(p.noteId).revision !== p.baseRevision)) {
      notify("Main changed. Update this branch and review it again.");
      return;
    }
    changes.forEach((p) => {
      const n = note(p.noteId);
      if (n) {
        n.title = p.title;
        n.body = p.body;
        n.source = p.source;
        n.revision++;
        n.verified = "5 Nov";
        n.stale = false;
        n.author = "Alex Morgan";
        n.initials = "AM";
      } else {
        s.notes.push({
          id: p.noteId,
          title: p.title,
          short: p.title,
          body: p.body,
          source: p.source,
          kind: p.kind,
          project: p.project,
          revision: 1,
          verified: "5 Nov",
          author: "Alex Morgan",
          initials: "AM",
          links: [],
          why: p.reason,
          commit: "Local concept revision",
          evidence: "Source reference: " + p.source,
        });
      }
      p.status = "merged";
    });
    b.status = "merged";
    event(
      "Alex merged " + name + ".",
      b.project +
        " · " +
        changes.length +
        " reviewed " +
        (changes.length === 1 ? "change" : "changes") +
        " joined main",
      "merge",
    );
    s.branch = "main";
    render();
    save();
    notify("Merged. Main now includes the reviewed changes. ✳");
  }
  document.addEventListener("click", (e) => {
    const target = e.target.closest("button,a,article[data-open-note]");
    if (!target) return;
    const d = target.dataset;
    if (d.close) {
      $("#" + d.close).close();
      return;
    }
    if (d.pin) {
      const i = s.pins.indexOf(d.pin);
      if (i < 0) s.pins.push(d.pin);
      else s.pins.splice(i, 1);
      const wasOpen = $("#note-dialog").open;
      render();
      if (wasOpen) openNote(d.pin);
      save();
      return;
    }
    if (d.view) {
      if (d.view === "memory") {
        s.project = "all";
        s.branch = "main";
        s.kind = "All";
        s.search = "";
        $("#global-search").value = "";
      }
      go(d.view);
      return;
    }
    if (d.project) {
      s.project = d.project;
      s.branch = "main";
      s.kind = "All";
      go("memory");
      return;
    }
    if (d.layout) {
      s.layout = d.layout;
      render();
      save();
      return;
    }
    if (d.kindFilter) {
      s.kind = d.kindFilter;
      render();
      save();
      return;
    }
    if (d.openNote) {
      openNote(d.openNote);
      return;
    }
    if (d.review) {
      $("#note-dialog").close();
      go("reviews", d.review);
      return;
    }
    if (d.openBranch) {
      go("branches", d.openBranch);
      return;
    }
    if ("backReviews" in d) {
      go("reviews");
      return;
    }
    if ("backBranches" in d) {
      go("branches");
      return;
    }
    if (d.reviewFilter) {
      s.reviewFilter = d.reviewFilter;
      s.selectedReview = "";
      render();
      save();
      return;
    }
    if ("compose" in d) {
      composer("", d.composeBranch || "");
      return;
    }
    if (d.editNote) {
      composer(d.editNote);
      return;
    }
    if (d.copyNote) {
      copyNote(d.copyNote);
      return;
    }
    if ("newBranch" in d) {
      $("#branch-form").reset();
      $("#branch-error").textContent = "";
      $("#branch-dialog").showModal();
      return;
    }
    if (d.browseBranch) {
      s.branch = d.browseBranch;
      s.project = branch(d.browseBranch).project;
      s.kind = "All";
      s.search = "";
      $("#global-search").value = "";
      go("memory");
      return;
    }
    if ("browseMain" in d) {
      s.branch = "main";
      s.project = "all";
      s.kind = "All";
      go("memory");
      return;
    }
    if ("localHowto" in d) {
      $("#local-howto-dialog").showModal();
      return;
    }
    if ("refreshLocal" in d) {
      loadLocal(true);
      return;
    }
    if (d.copyPath) {
      navigator.clipboard?.writeText(d.copyPath).then(
        () => notify("Path copied: " + d.copyPath),
        () => notify(d.copyPath),
      );
      return;
    }
    if (d.merge) {
      merge(d.merge);
      return;
    }
    if (d.rebase) {
      branchChanges(d.rebase).forEach((p) => {
        const n = note(p.noteId);
        if (n) {
          p.baseRevision = n.revision;
          p.baseBody = n.body;
        }
        if (p.status !== "merged") p.status = "pending";
      });
      event("Alex updated " + d.rebase + " from main.", "Drafts kept; review approvals cleared.", "branch");
      render();
      save();
      notify("Branch updated. Review the draft against its new base.");
      return;
    }
    if ("invite" in d) {
      $("#invite-form").reset();
      $("#invite-dialog").showModal();
      return;
    }
    if ("clearFilters" in d) {
      s.kind = "All";
      s.project = "all";
      s.search = "";
      $("#global-search").value = "";
      render();
      save();
      return;
    }
    if ("clearSearch" in d) {
      s.search = "";
      $("#global-search").value = "";
      render();
      save();
      return;
    }
    if ("export" in d) {
      const text =
        "# " + wsName() + " shared memory\n\n" + (LOCAL ? "Exported from Bakery (local).\n\n" : "Sample workspace export from Bakery.\n\n") +
        s.notes.map((n) => markdown(n)).join("\n\n---\n\n");
      const url = URL.createObjectURL(new Blob([text], { type: "text/markdown;charset=utf-8" }));
      const a = document.createElement("a");
      a.href = url;
      a.download = wsName().toLowerCase().replace(/[^a-z0-9]+/g, "-") + "-memory.md";
      a.click();
      setTimeout(() => URL.revokeObjectURL(url), 1000);
      notify("Sample shared memory packed as Markdown.");
      return;
    }
  });
  $("#global-search").value = s.search;
  $("#global-search").addEventListener("input", (e) => {
    s.search = e.target.value;
    s.view = "memory";
    s.selectedReview = "";
    s.selectedBranch = "";
    history.replaceState(null, "", "#memory");
    render();
    save();
  });
  document.addEventListener("change", (e) => {
    if (e.target.id === "memory-project") {
      s.project = e.target.value;
      render();
      save();
    }
    if (e.target.id === "memory-branch") {
      s.branch = e.target.value;
      s.project = s.branch === "main" ? "all" : branch(s.branch).project;
      s.kind = "All";
      render();
      save();
    }
  });
  $("#crumb-form").addEventListener("submit", (e) => {
    e.preventDefault();
    const data = new FormData(e.target),
      name = data.get("branch").trim(),
      project = data.get("project"),
      editId = data.get("editId"),
      existing = branch(name);
    if (!new RegExp("^[a-z0-9][a-z0-9/_-]*$").test(name)) {
      $("#crumb-error").textContent = "Use lowercase letters, numbers, slashes, dashes, or underscores.";
      return;
    }
    if (name === "main" || existing?.status === "merged") {
      $("#crumb-error").textContent = "Choose an open branch or a new name. Main stays separate.";
      return;
    }
    if (existing && existing.project !== project) {
      $("#crumb-error").textContent = "This branch belongs to another project. Choose another name.";
      return;
    }
    const title = data.get("title").trim(),
      body = data.get("body").trim(),
      source = data.get("source").trim();
    if (!title || !body || !source) {
      $("#crumb-error").textContent = "Add a claim, useful context, and its source.";
      return;
    }
    const id = editId || "note-" + Date.now().toString(36),
      current = note(id);
    let p = s.proposals.find((x) => x.noteId === id && x.branch === name && x.status !== "merged");
    if (!existing)
      s.branches.push({
        name,
        project,
        status: "open",
        author: "Alex Morgan",
        description: "A new change to project memory.",
      });
    const fields = {
      noteId: id,
      branch: name,
      title,
      body,
      source,
      kind: data.get("kind"),
      project,
      status: "pending",
      baseRevision: current?.revision || 0,
      baseBody: current?.body || "",
      author: "Alex Morgan",
      agent: "Bakery",
      reason: current
        ? "A proposed refinement to the existing context."
        : "A new piece of context worth keeping.",
    };
    if (p) Object.assign(p, fields);
    else {
      s.proposals.push({ ...fields, id: "p-" + Date.now().toString(36), comments: [] });
    }
    event(
      "Alex " + (editId ? "proposed an update to" : "left a new crumb:") + " " + title,
      project + " · " + name,
      "proposal",
    );
    $("#composer-dialog").close();
    s.branch = name;
    s.project = project;
    s.kind = "All";
    s.search = "";
    $("#global-search").value = "";
    go("memory");
    notify("Saved on the branch. Shared memory is unchanged until review and merge.");
  });
  $("#branch-form").addEventListener("submit", (e) => {
    e.preventDefault();
    const data = new FormData(e.target),
      name = data.get("name").trim();
    if (!new RegExp("^[a-z0-9][a-z0-9/_-]*$").test(name)) {
      $("#branch-error").textContent = "Use lowercase letters, numbers, slashes, dashes, or underscores.";
      return;
    }
    if (name === "main" || branch(name)) {
      $("#branch-error").textContent = "That branch name is already in use. Try a new one.";
      return;
    }
    s.branches.push({
      name,
      project: data.get("project"),
      status: "open",
      author: "Alex Morgan",
      description: "A fresh place to work out a useful change.",
    });
    event("Alex made a branch: " + name, data.get("project") + " · starts from main", "branch");
    $("#branch-dialog").close();
    go("branches", name);
    notify("A clean branch. Time to leave a good crumb.");
  });
  document.addEventListener("submit", (e) => {
    const form = e.target;
    if (!form.matches("[data-review-form]")) return;
    e.preventDefault();
    const p = proposal(form.dataset.reviewForm),
      action = e.submitter?.value,
      text = form.elements.comment.value.trim();
    if ((action === "comment" || action === "changes") && !text) {
      form.elements.comment.setCustomValidity("Leave a comment so the author knows what to change.");
      form.elements.comment.reportValidity();
      form.elements.comment.addEventListener("input", () => form.elements.comment.setCustomValidity(""), {
        once: true,
      });
      return;
    }
    if (text) p.comments.push({ author: "Alex", text, when: "Just now" });
    if (action === "approve") {
      p.status = "approved";
      event("Alex approved " + p.title, p.project + " · " + p.branch, "review");
    } else if (action === "changes") {
      p.status = "changes";
      event("Alex requested changes to " + p.title, p.project + " · " + p.branch, "review");
    } else event("Alex left a review comment.", p.project + " · " + p.branch, "comment");
    render();
    save();
    notify(
      action === "approve"
        ? "Approved. The branch can merge when all its changes are reviewed."
        : action === "changes"
          ? "Changes requested. Main keeps its existing version."
          : "Comment added to the review.",
    );
  });
  $("#invite-form").addEventListener("submit", (e) => {
    e.preventDefault();
    const data = new FormData(e.target),
      email = data.get("email").trim();
    if (s.members.some((m) => m.email.toLowerCase() === email.toLowerCase())) {
      notify("That teammate is already at the table.");
      return;
    }
    s.members.push({
      name: email.split("@")[0],
      email,
      role: data.get("role") + " · invited",
      initials: email.slice(0, 2).toUpperCase(),
    });
    $("#invite-dialog").close();
    event("Alex prepared a teammate invitation.", "Preview only · no email sent", "team");
    render();
    save();
    notify("Sample invitation added. No email was sent.");
  });
  $("#reset-demo").addEventListener("click", () => $("#confirm-dialog").showModal());
  $("#confirm-reset").addEventListener("click", () => {
    s = seed();
    try {
      localStorage.removeItem(KEY);
    } catch {}
    document.querySelectorAll("dialog[open]").forEach((d) => d.close());
    $("#global-search").value = "";
    history.replaceState(null, "", "#memory");
    render();
    save();
    notify("Fresh paper. The sample workspace has been reset.");
  });
  document.addEventListener("keydown", (e) => {
    if (
      e.key === "/" &&
      !["INPUT", "TEXTAREA", "SELECT"].includes(document.activeElement.tagName) &&
      !$("dialog[open]")
    ) {
      e.preventDefault();
      $("#global-search").focus();
    }
  });
  document.querySelectorAll("dialog").forEach((dialog) =>
    dialog.addEventListener("click", (e) => {
      if (e.target === dialog) {
        const r = dialog.getBoundingClientRect();
        if (e.clientX < r.left || e.clientX > r.right || e.clientY < r.top || e.clientY > r.bottom)
          dialog.close();
      }
    }),
  );
  ["mouseover", "focusin"].forEach((type) =>
    document.addEventListener(type, (e) => {
      const node = e.target.closest?.("[data-node]");
      if (node) traceWeb(node.dataset.node);
    }),
  );
  ["mouseout", "focusout"].forEach((type) =>
    document.addEventListener(type, (e) => {
      const node = e.target.closest?.("[data-node]");
      if (node && !node.contains(e.relatedTarget)) traceWeb("");
    }),
  );
  window.addEventListener("popstate", route);
  route();
  loadLocal(false);
  setInterval(() => {
    if (LOCAL && document.visibilityState === "visible") loadLocal(false);
  }, 4000);
  document.addEventListener("visibilitychange", () => {
    if (LOCAL && document.visibilityState === "visible") loadLocal(false);
  });
})();
