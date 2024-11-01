# Tareas
## Prioritarias
### Global
Implementar errores de GRPC.
Implenetar gestion de errores global.
### Implant
solucionar busy_waiting
### Client
revisar si merece la pena handleStatusCodes(s *status.Status) int{
### Server
Arreglar  metodo GetImplants
## Secundarias Importantes
### Implant
En caso de afrontar el busy_waiting con sleep, generar implant con un tiempo especifico de sleep
Si un implant muere por cualquier causa hay que quitarlo de la lista 
## Secundarias Prescindibles
### Client
Funcion de recarga para actualizar los implants de tipo tea.Cmd, implementar spinners de 
carga o algo similar mientras espera la respuesta del servidor
Se puede implementar la reverse shell con tea.Exec
Mensaje cuando no hay implants conectados
