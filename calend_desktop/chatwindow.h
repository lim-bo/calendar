#ifndef CHATWINDOW_H
#define CHATWINDOW_H

#include <QDialog>
#include "client.h"
#include "cfg.h"
#include <QMessageBox>
#include "message_entry.h"
#include <QLayout>
namespace Ui {
class ChatWindow;
}

class ChatWindow : public QDialog
{
    Q_OBJECT

public:
    explicit ChatWindow(QString eventID, QString uid, QWidget *parent = nullptr);
    ~ChatWindow();

private slots:
    void on_sendButton_clicked();

private:
    Ui::ChatWindow *ui;
    QString eventID;
    QString uid;
    Client cli;
    void loadMessages();
};

#endif // CHATWINDOW_H
