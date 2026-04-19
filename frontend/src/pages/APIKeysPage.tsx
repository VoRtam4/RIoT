import { useState } from "react";
import APIKeyList from "../modules/apiKeys/components/APIKeyList";
import APIKeyForm from "../modules/apiKeys/components/APIKeyForm";
import APIKeyCreatedBox from "../modules/apiKeys/components/APIKeyCreatedBox";
import { useApiKeys } from "../modules/apiKeys/hooks/useApiKeys";
import { useApiKeyMutations } from "../modules/apiKeys/hooks/useApiKeyMutations";
import toast from "react-hot-toast";
import { useNavigate } from "react-router-dom";

export default function APIKeysPage() {
  const navigate = useNavigate();
  const { apiKeys, loading, refetch } = useApiKeys();
  const { createApiKey, updateApiKey, deleteApiKey } = useApiKeyMutations();

  const [mode, setMode] = useState<"empty" | "create" | "detail" | "created">(
    "empty",
  );
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [createdKey, setCreatedKey] = useState<string | null>(null);

  const selected = apiKeys.find((k: any) => k.id === selectedId);

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

            <APIKeyList
              apiKeys={apiKeys}
              selectedId={selectedId}
              loading={loading}
              onSelect={(id) => {
                if (selectedId === id) {
                  setSelectedId(null);
                  setMode("empty");
                  return;
                }
                setSelectedId(id);
                setMode("detail");
              }}
            />
          </div>
        </div>

        <div className="col-9">
          <div className="panel h-100 p-3 d-flex flex-column">
            {mode === "empty" && (
              <div className="d-flex justify-content-center align-items-center h-100">
                <button
                  className="btn btn-primary"
                  onClick={() => setMode("create")}
                >
                  New key
                </button>
              </div>
            )}

            {mode === "create" && (
              <APIKeyForm
                mode="create"
                onSubmit={async (data) => {
                  try {
                    const res = await createApiKey({
                      variables: { input: data },
                    });
                    setCreatedKey(res.data?.createAPIKey ?? null);
                    setMode("created");
                    await refetch();
                    toast.success("Created");
                  } catch {
                    toast.error("Error");
                  }
                }}
              />
            )}

            {mode === "detail" && selected && (
              <APIKeyForm
                mode="detail"
                apiKey={selected}
                onSubmit={async (data) => {
                  try {
                    await updateApiKey({
                      variables: { id: selected.id, input: data },
                    });
                    await refetch();
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
                    setSelectedId(null);
                    await refetch();
                    toast.success("Deleted");
                  } catch {
                    toast.error("Error");
                  }
                }}
              />
            )}

            {mode === "created" && createdKey && (
              <APIKeyCreatedBox value={createdKey} />
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
