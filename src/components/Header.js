import React, { useState } from "react";
import "./Header.css";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
} from "@mui/material";

function Header() {
  const [loginDialogOpen, setLoginDialogOpen] = useState(false);
  const [registerDialogOpen, setRegisterDialogOpen] = useState(false);

  //funkcija koje ce otvoriti dijalog prilikom klika na login button iz headera
  const openLoginDialog = () => {
    setLoginDialogOpen(true);
  };

  // TODO: ovdje na analogan nacin napravi funkciju za otvarnje dijaloga za registraciju

  // funkcije koje ce zatvoriti dijalog prilikom klika na login button iz headera
  const closeLoginDialog = () => {
    setLoginDialogOpen(false);
  };

  // TODO: ovdje na analogan nacin napravi funkciju za zatvaranje dijaloga za registraciju

  return (
    <>
      <div className="header">
        <div className="logo">NewBetting</div>
        <Button onClick={openLoginDialog}>Prijava</Button>
        {/* ovdje dolje u onClick ubaci funkciju kojom ce se otvoriti registracijski dijalog */}
        <Button onClick={null}>Registracija</Button>
      </div>

      {/* login dialog */}
      <div>
        <Dialog open={loginDialogOpen} onClose={closeLoginDialog}>
          <DialogTitle>Prijava</DialogTitle>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              id="name"
              label="Korisničko ime / email adresa"
              type="email"
              fullWidth
              variant="standard"
            />
            <TextField
              autoFocus
              margin="dense"
              id="password"
              label="Lozinka"
              type="password"
              fullWidth
              variant="standard"
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={closeLoginDialog}>Odustani</Button>
            {/* zasad "prijavi se" button nista ne radi, tu ces kasnije (kada dodjes do tog todo-a) dodati funkciju koja ce slati HTTP request na backend; zato zasad stavljamo funkciju na klik da bude null */}
            <Button onClick={null}>Prijavi se</Button>
          </DialogActions>
        </Dialog>
      </div>

      <div>
        {/* na analogan nacin napravi ovdje registration Dialog komponentu i neka u njemu budu sva polja koja su zadana u Basecamp todo-u */}
      </div>
    </>
  );
}

export default Header;
