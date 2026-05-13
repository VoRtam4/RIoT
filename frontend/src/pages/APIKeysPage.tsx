/**
 * @file APIKeysPage.tsx
 * @brief Stránka správy API klíčů, oprávnění, IP omezení a životnosti klíčů.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useState } from "react";
import APIKeySidebar from "../modules/apiKeys/components/APIKeySidebar";
import APIKeyForm from "../modules/apiKeys/components/APIKeyForm";
import APIKeyCreatedBox from "../modules/apiKeys/components/APIKeyCreatedBox";
import { useApiKeys } from "../modules/apiKeys/hooks/useApiKeys";
import { useCreateApiKey } from "../modules/apiKeys/hooks/useCreateApiKey";
import { useUpdateApiKey } from "../modules/apiKeys/hooks/useUpdateApiKey";
import { useDeleteApiKey } from "../modules/apiKeys/hooks/useDeleteApiKey";
import toast from "react-hot-toast";
import { useNavigate } from "react-router-dom";
import { usePageState } from "../app/navigation/usePageState";
import { apiKeysPageStateCodec } from "../modules/apiKeys/state/apiKeysPageState";

export default function APIKeysPage() {
  const navigate = useNavigate();

  const { apiKeys, loading } = useApiKeys();
  const { createApiKey } = useCreateApiKey();
  const { updateApiKey } = useUpdateApiKey();
  const { deleteApiKey } = useDeleteApiKey();
  const { query, setPageState } = usePageState(apiKeysPageStateCodec);
  const [mode, setMode] = useState<"empty" | "create" | "created">("empty");
  const [createdKey, setCreatedKey] = useState<string | null>(null);
  const selected = apiKeys.find(
    (k: any) => String(k.id) === String(query.selected),
  );
  const effectiveMode =
    mode === "create" || mode === "created"
      ? mode
      : selected
        ? "detail"
        : "empty";

  return (
    <div className="container mt-3">
      <div className="row g-3 h-100" style={{ height: "calc(100vh - 80px)" }}>
        <div className="col-3">
          <div
            className="panel p-3 d-flex flex-column"
            style={{
              position: "sticky",
              top: 16,
              height: "calc(100vh - 110px)",
            }}
          >
            <button
              className="btn btn-outline-light w-100 mb-3"
              onClick={() => navigate("/api-keys/docs")}
            >
              API documentation
            </button>

            <APIKeySidebar
              apiKeys={apiKeys}
              selectedId={query.selected}
              loading={loading}
              search={query.q}
              filter={query.filter}
              sort={query.sort}
              onSearchChange={(q) =>
                setPageState(
                  {
                    query: {
                      ...query,
                      q,
                    },
                    entry: null,
                  },
                  { replace: true },
                )
              }
              onFilterChange={(filter) =>
                setPageState(
                  {
                    query: {
                      ...query,
                      filter,
                    },
                    entry: null,
                  },
                  { replace: true },
                )
              }
              onSortChange={(sort) =>
                setPageState(
                  {
                    query: {
                      ...query,
                      sort,
                    },
                    entry: null,
                  },
                  { replace: true },
                )
              }
              onSelect={(id) => {
                setMode("empty");

                if (String(query.selected) === String(id)) {
                  setPageState(
                    {
                      query: {
                        ...query,
                        selected: null,
                      },
                      entry: null,
                    },
                    { replace: true },
                  );
                  return;
                }

                setPageState(
                  {
                    query: {
                      ...query,
                      selected: String(id),
                    },
                    entry: null,
                  },
                  { replace: true },
                );
              }}
            />
          </div>
        </div>

        <div className="col-9">
          <div className="panel h-100 p-3 d-flex flex-column">
            {effectiveMode === "empty" && (
              <div className="d-flex flex-column h-100">
                <div className="flex-grow-1" />
                <button
                  className="btn btn-primary"
                  style={{ alignSelf: "flex-end" }}
                  onClick={() => {
                    setCreatedKey(null);
                    setMode("create");
                  }}
                >
                  New key
                </button>
              </div>
            )}

            {effectiveMode === "create" && (
              <APIKeyForm
                mode="create"
                onSubmit={async (data) => {
                  try {
                    const res = await createApiKey({
                      variables: { input: data },
                    });

                    setCreatedKey(res.data?.createAPIKey ?? null);
                    setMode("created");

                    toast.success("Created");
                  } catch {
                    toast.error("Error");
                  }
                }}
              />
            )}

            {effectiveMode === "detail" && selected && (
              <APIKeyForm
                mode="detail"
                apiKey={selected}
                onSubmit={async (data) => {
                  try {
                    await updateApiKey({
                      variables: { id: selected.id, input: data },
                    });

                    toast.success("Saved");
                  } catch {
                    toast.error("Error");
                  }
                }}
                onDelete={async () => {
                  try {
                    await deleteApiKey({
                      variables: { id: selected.id },
                    });

                    setMode("empty");
                    setPageState(
                      {
                        query: {
                          ...query,
                          selected: null,
                        },
                        entry: null,
                      },
                      { replace: true },
                    );

                    toast.success("Deleted");
                  } catch {
                    toast.error("Error");
                  }
                }}
              />
            )}

            {effectiveMode === "created" && createdKey && (
              <APIKeyCreatedBox
                value={createdKey}
                onNew={() => {
                  setCreatedKey(null);
                  setMode("create");
                }}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
