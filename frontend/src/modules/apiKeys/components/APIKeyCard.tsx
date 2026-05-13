/**
 * @file APIKeyCard.tsx
 * @brief Karta API klíče se stavem, oprávněními a základními akcemi.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";

import { colors } from "../../../theme/colors";

type Props = {
  apiKey: any;
  selected: boolean;
  onClick: () => void;
};

export default function APIKeyCard({ apiKey, selected, onClick }: Props) {
  const isActive =
    !apiKey.revoked &&
    (!apiKey.expiresAt || new Date(apiKey.expiresAt) > new Date());

  return (
    <Card
      sx={{
        flexShrink: 0,
        height: "100%",
        overflow: "hidden",
        backgroundColor: selected
          ? "var(--border-input)"
          : "var(--bg-card)",
        "&:hover": {
          backgroundColor: selected
          ? "var(--border-input)"
          : "var(--border)",
        },
      }}
    >
      <CardActionArea onClick={onClick} sx={{ height: "100%" }}>
        <CardContent
          sx={{
            py: 1.5,
            height: "100%",
            overflow: "hidden",
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
          }}
        >
          {/* LABEL */}
          <Typography
            variant="body1"
            sx={{
              display: "-webkit-box",
              overflow: "hidden",
              textOverflow: "ellipsis",
              WebkitBoxOrient: "vertical",
              WebkitLineClamp: 2,
            }}
          >
            {apiKey.label || "\u00A0"}
          </Typography>

          {/* STATUS */}
          <Box
            sx={{
              display: "flex",
              alignItems: "center",
              gap: 1,
              mt: 0.5,
            }}
          >
            <Box
              sx={{
                width: 8,
                height: 8,
                borderRadius: "50%",
                backgroundColor: isActive
                  ? colors.success
                  : colors.error,
              }}
            />

            <Typography variant="caption" sx={{ opacity: 0.8 }}>
              {isActive ? "Active" : "Inactive"}
            </Typography>
          </Box>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}
