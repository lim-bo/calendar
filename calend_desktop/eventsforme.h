#ifndef EVENTSFORME_H
#define EVENTSFORME_H

#include <QWidget>
#include <QVBoxLayout>
#include <QLabel>
#include "event_entry.h"
#include "client.h"
#include "eventData.h"

namespace Ui {
class eventsforme;
}

class eventsforme : public QWidget
{
    Q_OBJECT

public:
    explicit eventsforme(QString uid, QWidget *parent = nullptr);
    ~eventsforme();

private slots:
    void on_checkBox_created_stateChanged(int);

    void on_checkBox_accepted_stateChanged(int);

    void on_checkBox_not_accepted_stateChanged(int);

private:
    Ui::eventsforme *ui;
    QString uid;
    Client client;

    void loadEvents();

};

#endif // EVENTSFORME_H
