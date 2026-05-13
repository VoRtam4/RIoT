/**
 * @file main.tsx
 * @brief Vstupní bod frontendové aplikace platformy RIoT.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @defgroup riot_frontend Frontend
 * @ingroup riot
 * @see ../README.md
 */
import React from "react";
import ReactDOM from "react-dom/client";
import { ApolloProvider } from "@apollo/client/react";
import { RouterProvider } from "react-router-dom";
import { apolloClient } from "./app/apollo";
import { router } from "./app/router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "react-hot-toast";
import { ThemeProvider } from "@mui/material";
import { muiTheme } from "./theme/muiTheme";
import ThemeInitializer from "./theme/ThemeInitializer";
import TimeSeriesExportManager from "./modules/timeSeries/components/TimeSeriesExportManager";
import "@fortawesome/fontawesome-free/css/all.min.css";
import "bootstrap/dist/css/bootstrap.min.css";
import "./styles/theme.css";

const queryClient = new QueryClient();

ReactDOM.createRoot(document.getElementById("root")!).render(
  <ThemeProvider theme={muiTheme}>
    <React.StrictMode>
      <ThemeInitializer />
      <ApolloProvider client={apolloClient}>
        <QueryClientProvider client={queryClient}>
          <>
            <RouterProvider router={router} />
            <TimeSeriesExportManager />
            <Toaster
              position="bottom-center"
              toastOptions={{
                duration: 3000,
                style: {
                  background: "var(--bg-card)",
                  color: "var(--text-main)",
                  border: "1px solid var(--border)",
                  borderRadius: "10px",
                  padding: "10px 14px",
                  fontSize: "14px",
                },
                success: {
                  icon: (
                    <i
                      className="fa-solid fa-check"
                      style={{ color: "var(--success)" }}
                    />
                  ),
                },
                error: {
                  icon: (
                    <i
                      className="fa-solid fa-xmark"
                      style={{ color: "var(--error)" }}
                    />
                  ),
                },
              }}
            />
          </>
        </QueryClientProvider>
      </ApolloProvider>
    </React.StrictMode>
  </ThemeProvider>,
);
